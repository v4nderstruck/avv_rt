<script>
  import { Block, List, ListItem } from "konsta/svelte";
  import levehenstein from "../utils/levehenstein";
  import { Busstops } from "../utils/busstops";

  export let searchInput;

  let possibleStations = Busstops;
  let displayStations = [];

  const compareWithInput = (input, station_0, station_1) => {
    const distance_a = levehenstein(input, station_0);
    const distance_b = levehenstein(input, station_1);
    return distance_a - distance_b;
  };

  function filterPossibleStations(input) {
    const start = Date.now();
    possibleStations = possibleStations.sort((s_0, s_1) =>
      compareWithInput(input, s_0, s_1)
    );
    console.log(input, "in ", Date.now() - start, "ms");
    displayStations = possibleStations.slice(0, 5);
    displayStations.forEach((v) => console.log(v, levehenstein(input, v)));
  }

  $: filterPossibleStations(searchInput);
</script>

<Block strong inset class="z-50 absolute mb-0 mt-4 pt-0 pb-0 fill-width">
  <List class="mt-0 mb-0">
    {#each displayStations as station (station)}
      <ListItem title={station} />
    {/each}
  </List>
</Block>
