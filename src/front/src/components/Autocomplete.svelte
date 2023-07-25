<script>
  import { Block, List, ListItem } from "konsta/svelte";

  export let searchInput;

  let displayStations = [];

  const updateStationList = async (input) => {
    const res = await fetch(`/api/autocomplete?q=${input}`);
    const data = await res.json();
    if (data && data.scores) {
      displayStations = data.scores;
    }
  };

  $: updateStationList(searchInput);
</script>

<Block strong inset class="z-50 absolute mb-0 mt-4 pt-0 pb-0 fill-width">
  <List class="mt-0 mb-0">
    {#each displayStations as station (station)}
      <ListItem title={station.name} link href={`/stop/${station.id}`} />
    {/each}
  </List>
</Block>
