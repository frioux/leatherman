Lists coffee from sweetmarias as JSON.

Run lists all of the available [Sweet Maria's](https://www.sweetmarias.com/) coffees
as json documents per line.  Here's how you might see the top ten coffees by
rating:

```bash
$ sm-list | jq -r '[.Score, .Title, .URL ] | @tsv' | sort -n | tail -10
87.7    Colombia EA Café Quindio Decaf  https://www.sweetmarias.com/colombia-ea-cafe-quindio-decaf-6728.html
87.7    Guatemala Antigua Pulcal Inteligente    https://www.sweetmarias.com/guatemala-antigua-pulcal-inteligente-6604.html
88.2    El Salvador Santa Ana Pacamara AAA Lot 2        https://www.sweetmarias.com/el-salvador-santa-ana-pacamara-aaa-lot-2-6653.html
88.7    Burundi Monge Murambi Hill      https://www.sweetmarias.com/burundi-monge-murambi-hill-6643.html
89      Ethiopia Uraga Tebe Haro Wato   https://www.sweetmarias.com/ethiopia-uraga-tebe-haro-wato-6725-.html
89.4    Burundi Kazoza N'Ikawa Coop     https://www.sweetmarias.com/burundi-kazoza-nikawa-station-6639.html
89.6    Ethiopia Nano Challa Cooperative        https://www.sweetmarias.com/ethiopa-nano-challa-cooperative-6726.html
89.8    Ethiopia Dry Process Yirga Cheffe Aricha        https://www.sweetmarias.com/ethiopia-dry-process-yirga-cheffe-aricha-6680.html
90.1    Ethiopia Dry Process Tarekech Werasa    https://www.sweetmarias.com/ethiopia-dry-process-tarekech-werasa-6701.html
90.2    Ethiopia Organic Shakiso Kayon Mountain Farm    https://www.sweetmarias.com/ethiopia-organic-shakiso-kayon-mountain-farm-6681.html
```
