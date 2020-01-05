package sweetmarias

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"net/http/httptest"

	"github.com/frioux/leatherman/internal/testutil"
)

func TestLoadCoffee(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		f, err := os.Open("./testdata/sm.html")
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(w, f)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	c, err := LoadCoffee(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("Failed to LoadCoffee: %s", err)
	}

	testutil.Equal(t, c, Coffee{
		Title:    "Papua New Guinea Honey Nebilyer Estate",
		Overview: "Honey process seems to bring out fruited notes like cranberry, raisin, red grape, and underscored by molasses-like sweetness. This PNG boasts body, and with mild acidity, is great espresso too. City+ to Full City+. Good for espresso.",
		Score:    86.6,
		URL:      ts.URL,

		FarmNotes:    "This coffee comes to us by way of the Kuta coffee mill in the Waghi District of Papua New Guinea. The coffee processed at the mill are from smaller coffee plantations in the area situated at just under 1600 meters above sea level on the low end. It's a honey processed coffee, meaning the coffee cherry and much of the fruit are stripped from the seed using depulping machinery, and then the seed still covered in sticky mucilage is laid to dry with any remaining fruit still intact. This tends to result in bigger body, softer acidity, and often a fruited cup. The physical grade of this coffee is impressive, and I couldn't find a single full quaker bean in the few hundred grams of coffee I roasted.",
		CuppingNotes: "A honey process batch from the same mill that brought us Kuta Waghi (we still have some available). In fact, tasting these two coffees side by side, you get an idea of the role processing plays in a coffee's final cup profile. This lot is much more fruited than the Peaberry, and the aroma displays a sweet blend of honey and molasses sweetness, and fruited smells like dark berry pulp drifting up in the steam. The cup displays a nice fruited profile as well when roasted to City+, and a sort of unrefined sweetness at the core helps to highlight fruited nuance. The cooling coffee reveals glimpses of red grape, raisin, and a dark cranberry note in the finish. Full City roasts are also highlighted by berry characteristics, along with a rustic dark chocolate base note. Acidity is quite mild (typical of honey process coffee) across the roast spectrum, and Full City roasts produce inky espresso shots with layers deep cocoa roast tone interspersed with a sweet cranberry hint.",

		AdditionalAttributes: map[string]string{
			"Appearance":               ".8 d per 300 grams, 15 - 19 Screen",
			"Arrival date":             "November 2018 Arrival",
			"Bag size":                 "60 KG",
			"Cultivar Detail":          "Arusha, Bourbon, Typica",
			"Drying Method":            "Patio Sun-dried",
			"Grade":                    "A/X",
			"Lot size":                 "32",
			"Packaging":                "GrainPro liner",
			"Processing":               "Honey Process",
			"Recommended for Espresso": "Yes",
			"Region":                   "Waghi Valley",
			"Roast Recommendations":    "City+ to Full City+",
		},
	}, "wrong coffee")
}
