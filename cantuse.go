package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type Data struct {
	// {
	//     "chrome": {
	//         "browser": "Chrome",
	//         "abbr": "Chr.",
	//         "usage_global": {
	//             "10": 0.004534,
	//             "11": 0.004464,
	//             "12": 0.010424,
	//            ...
	//         }
	//      },
	//      {
	//          ...
	//      },
	// }
	Agents map[string]struct {
		Browser     string             `json:"browser"`
		Abbr        string             `json:"abbr"`
		UsageGlobal map[string]float32 `json:"usage_global"`
	} `json:"agents"`

	// {
	//     "let": {
	//         "title": "let",
	//         "description": "Declares a variable with block level scope",
	//         "spec": "https://www.ecma-international.org/ecma-262/6.0/#sec-let-and-const-declarations",
	//         "links": [
	//             {
	//                 "url": "http://generatedcontent.org/post/54444832868/variables-and-constants-in-es6",
	//                 "title": "Variables and Constants in ES6"
	//             },
	//             {
	//                 "url": "https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/let",
	//                 "title": "MDN Web Docs - let"
	//             }
	//         ],
	//         "stats": {
	//            "chrome": {
	//                "10": "n",
	//                "11": "n",
	//                ...
	//                "22": "n d #2",
	//                "23": "n d #2",
	//                ...
	//                "75": "y",
	//                "76": "y",
	//                ...
	//            },
	//            {
	//                ...
	//            },
	//         }
	//      }
	Data map[string]struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Spec        string  `json:"spec"`
		UsagePercY  float32 `json:"usage_perc_y"`
		UsagePercA  float32 `json:"usage_perc_a"`

		Links []struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"links"`
		Stats map[string]map[string]string `json:"stats"`
	} `json:"data"`
}

func main() {
	var (
		help, untracked bool
		ignore          string
	)
	flag.BoolVar(&untracked, "untracked", false, `Count "untracked" browsers as supported.`)
	flag.StringVar(&ignore, "ignore", "", "List of browsers to ignore")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	ignoreList := strings.Split(strings.ToLower(ignore), ",")

	d, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	var data Data
	err = json.Unmarshal(d, &data)
	if err != nil {
		log.Fatal(err)
	}

	ut := 100.0 - data.Data["css-sel2"].UsagePercY
	if !untracked {
		fmt.Printf("Untracked: %.1f%%; support can never be higher than %.1f%%.\n\n",
			ut, data.Data["css-sel2"].UsagePercY)
	} else {
		fmt.Printf("Untracked: %.1f%%\n\n", ut)
	}

	var feats []string
	for f := range data.Data {
		feats = append(feats, f)
	}
	sort.Strings(feats)

	for _, feat := range feats {
		desc := data.Data[feat]

		type w struct {
			s string
			n float32
		}
		var wontwork []w
		var total, partial float32
	nextBrowser:
		for browser, versions := range desc.Stats {
			for version, supported := range versions {
				usage := data.Agents[browser].UsageGlobal[version]
				if strings.HasPrefix(supported, "y") {
					total += usage
				} else if strings.HasPrefix(supported, "a") {
					partial += usage
				} else if usage > 0.050 {
					b := data.Agents[browser].Browser + " " + version
					for _, ig := range ignoreList {
						if ig == strings.ToLower(b) {
							total += usage
							continue nextBrowser
						}
					}

					wontwork = append(wontwork, w{
						fmt.Sprintf("%s %s %.1f%%", b, strings.Repeat(" ", 30-len(b)), usage),
						usage})
				}
			}
		}
		sort.Slice(wontwork, func(i, j int) bool { return wontwork[i].n > wontwork[j].n })

		if untracked {
			total += ut
		}
		sup := 100.0 - total - partial
		if sup <= 0 {
			fmt.Printf("%q should work for all visitors", feat)
		} else if sup < 0.05 {
			fmt.Printf("%q should work for almost all visitors (unsupported: %.3f%%)", feat, sup)
		} else {
			fmt.Printf("%q won't work for %.1f%% of visitors", feat, sup)
		}
		if partial > 0 {
			fmt.Printf(" (partial: %.1f%%)", partial)
		}
		fmt.Println("")
		for i, w := range wontwork {
			if i >= 10 {
				fmt.Printf("\t…%d more…\n", len(wontwork)-10)
				break
			}
			fmt.Printf("\t%s\n", w.s)
		}
		fmt.Println("")
	}
}
