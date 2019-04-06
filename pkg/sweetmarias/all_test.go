package sweetmarias

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllCoffees(t *testing.T) {
	f, err := os.Open("./testdata/all.html")
	if err != nil {
		t.Fatalf("Couldn't open all.html: %s", err)
	}

	c, err := allCoffees(f)
	if err != nil {
		t.Fatalf("Couldn't load allCoffees: %s", err)
	}
	assert.Equal(t, []string{
		"https://www.sweetmarias.com/sweet-maria-s-new-classic-espresso.html",
		"https://www.sweetmarias.com/rwanda-nyamasheke-cyato-2018.html",
		"https://www.sweetmarias.com/costa-rica-helsar-cascara-coffee-fruit-tea-2018.html",
		"https://www.sweetmarias.com/congo-kivu-minova-cpu.html",
		"https://www.sweetmarias.com/flores-wolo-bobo.html",
		"https://www.sweetmarias.com/kenya-kiambu-japem-c-grade-5878.html",
		"https://www.sweetmarias.com/guatemala-xinabajul-productores-de-cuilco-5986.html",
		"https://www.sweetmarias.com/rwanda-nyungwe-swp-decaf.html",
		"https://www.sweetmarias.com/peru-huabal-la-esperanza-6061.html",
		"https://www.sweetmarias.com/sweet-maria-s-ethiopiques-blend.html",
		"https://www.sweetmarias.com/timor-leste-goulala-village.html",
		"https://www.sweetmarias.com/ethiopia-agaro-nano-challa-coop-2018.html",
		"https://www.sweetmarias.com/brazil-don-jose-maria-yellow-catuai.html",
		"https://www.sweetmarias.com/sweet-marias-moka-kadir-blend.html",
		"https://www.sweetmarias.com/colombia-buesaco-cristian-adarme-6035.html",
		"https://www.sweetmarias.com/roasted-espresso-sweet-maria-s-espresso-workshop-44.html",
		"https://www.sweetmarias.com/roasted-coffee-burundi-kayanza-n-ikawa-station.html",
		"https://www.sweetmarias.com/guatemala-huehuetenango-xinabajul-producers-6000.html",
		"https://www.sweetmarias.com/brazil-pulp-natural-fazenda-santa-ines-6068.html",
		"https://www.sweetmarias.com/rstd-subs-1050.html",
		"https://www.sweetmarias.com/sumatra-wet-process-gunung-tujuh-2018.html",
		"https://www.sweetmarias.com/sweet-maria-s-decaf-espresso-donkey-blend.html",
		"https://www.sweetmarias.com/colombia-honey-process-aponte-community-6034.html",
		"https://www.sweetmarias.com/rwanda-nyamasheke-gitwe.html",
		"https://www.sweetmarias.com/kenya-kiambu-mandela-estate-ab-5886.html",
		"https://www.sweetmarias.com/colombia-inza-rio-paez-6024.html",
		"https://www.sweetmarias.com/green-coffee-sample-set-regular.html",
		"https://www.sweetmarias.com/flores-wolo-wio.html",
		"https://www.sweetmarias.com/bali-wet-hulled-bangli-kintamani.html",
		"https://www.sweetmarias.com/catalog/product/view/id/13940/s/flores-poco-ranaka/category/4/",
		"https://www.sweetmarias.com/ethiopia-dry-process-dambi-uddo-site-6038.html",
		"https://www.sweetmarias.com/rwanda-gikongoro-robusta.html",
		"https://www.sweetmarias.com/timor-letefoho-poulala.html",
		"https://www.sweetmarias.com/sweet-maria-s-liquid-amber-espresso-blend.html",
		"https://www.sweetmarias.com/sweet-maria-s-altiplano-blend.html",
		"https://www.sweetmarias.com/ethiopia-agaro-sadi-loya-coop-2018.html",
		"https://www.sweetmarias.com/guatemala-patzun-finca-la-florida-5860.html",
		"https://www.sweetmarias.com/roasted-coffee-colombia-inza-veredas-vecinas.html",
		"https://www.sweetmarias.com/peru-jaen-granjeros-de-huabal-6063.html",
		"https://www.sweetmarias.com/burundi-kayanza-nemba-station-5957.html",
		"https://www.sweetmarias.com/el-salvador-la-esperanza-swp-decaf.html",
		"https://www.sweetmarias.com/green-coffee-sample-set-espresso.html",
		"https://www.sweetmarias.com/kenya-nyeri-aberdare-aa-5832.html",
		"https://www.sweetmarias.com/sumatra-giling-basah-kerinci.html",
		"https://www.sweetmarias.com/guatemala-dry-process-finca-rosma-5997.html",
		"https://www.sweetmarias.com/sweet-maria-s-espresso-monkey-blend.html",
		"https://www.sweetmarias.com/espresso-workshop-44-carga-larga-6071.html",
		"https://www.sweetmarias.com/burundi-kazoza-n-ikawa-station-5962.html",
		"https://www.sweetmarias.com/java-wet-hulled-frinsa-estate-6051.html",
		"https://www.sweetmarias.com/colombia-la-plata-raul-hector-6025.html",
		"https://www.sweetmarias.com/ethiopia-yirga-cheffe-kochore-boji-5796.html",
		"https://www.sweetmarias.com/nicaragua-buenos-aires-san-salvador-5803.html",
		"https://www.sweetmarias.com/burundi-kiganda-murambi-5966.html",
		"https://www.sweetmarias.com/guatemala-huehuetenango-la-libertad-lot-2-2018.html",
		"https://www.sweetmarias.com/ethiopia-guji-uraga-tome-station-5784.html",
		"https://www.sweetmarias.com/brazil-dry-process-pedra-branca-lot-1-6069.html",
		"https://www.sweetmarias.com/sumatra-jagong-jeget-swp-decaf.html",
		"https://www.sweetmarias.com/ethiopia-nansebo-tulu-golla-5779.html",
		"https://www.sweetmarias.com/rwanda-karongi-gitesi-2018.html",
		"https://www.sweetmarias.com/burundi-kayanza-kibingo-station-5956.html",
		"https://www.sweetmarias.com/rwanda-rulindo-tumba-station-lot-442.html",
		"https://www.sweetmarias.com/brazil-dry-process-fazenda-campos-altos.html",
		"https://www.sweetmarias.com/guatemala-proyecto-xinabajul-senor-aler-5874.html",
		"https://www.sweetmarias.com/guatemala-huehuetenango-xinabajul-swp-decaf-2018.html",
		"https://www.sweetmarias.com/colombia-honey-process-buenos-aires-gesha-6037.html",
		"https://www.sweetmarias.com/green-coffee-decaf-sample-set-4-pounds-1.html",
		"https://www.sweetmarias.com/papua-new-guinea-honey-process-nebilyer-estate.html",
	}, c)
}
