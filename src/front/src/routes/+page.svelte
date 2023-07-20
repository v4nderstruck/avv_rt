<script>
  import {
    Page,
    Navbar,
    List,
    ListItem,
    BlockTitle,
    ListInput,
  } from "konsta/svelte";
  import FaBus from "svelte-icons/fa/FaBus.svelte";
  import Autocomplete from "../components/Autocomplete.svelte";

  let stationSearch = "";

  const onStationSearchChanged = (e) => {
    stationSearch = e.target.value;
  };

  const onStationSearchCleared = () => {
    stationSearch = "";
  };
</script>

<Page>
  <Navbar title="🚌 AC Bus Departures" />
  <BlockTitle large>Bus Departures</BlockTitle>
  <div class="z-20 w-full">
    <List strongIos insetIos class="mt-8 mb-0 overflow-visible">
      <ListInput
        label="Find Your Bus stop"
        type="text"
        value={stationSearch}
        clearButton={stationSearch.length > 0}
        onInput={onStationSearchChanged}
        onClear={onStationSearchCleared}
      >
        <div class="w-8 h-8" slot="media">
          <FaBus />
        </div>
      </ListInput>
    </List>

    {#if stationSearch.length > 0}
      <Autocomplete searchInput={stationSearch} />
    {/if}
  </div>

  <BlockTitle large>Current Delays</BlockTitle>
  <List strongIos insetIos>
    <ListItem title="Delay" />
  </List>
</Page>

<style>
  .icon {
    width: 16px;
    height: 16px;
  }
</style>
