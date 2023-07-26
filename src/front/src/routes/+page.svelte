<script>
  import {
    Page,
    Navbar,
    List,
    ListItem,
    BlockTitle,
    ListInput,
    Block,
  } from "konsta/svelte";
  import FaRegSmile from "svelte-icons/fa/FaRegSmile.svelte";
  import FaBus from "svelte-icons/fa/FaBus.svelte";
  import Autocomplete from "../components/Autocomplete.svelte";
  import { onMount } from "svelte";
  import FaRegStar from "svelte-icons/fa/FaRegStar.svelte";

  let stationSearch = "";
  let timer;
  let localStorage = null;
  let watchList = {};

  const onStationSearchChanged = (e) => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      stationSearch = e.target.value;
    }, 250);
  };

  onMount(() => {
    localStorage = window.localStorage;
    if (localStorage && localStorage.getItem("stationWatchList")) {
      watchList = JSON.parse(localStorage.getItem("stationWatchList"));
    }
  });

  $: showWatchList = () => {
    console.log("WatchList", watchList);
    return Object.keys(watchList).length > 0;
  };

  const onStationSearchCleared = () => {
    stationSearch = "";
  };
</script>

<div>
  <BlockTitle large>Bus Departures</BlockTitle>
  <div class="z-20 w-full">
    <List strongIos insetIos class="mt-8 mb-0 overflow-visible">
      <ListInput
        label="Enter Bus stop"
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

  {#if showWatchList()}
    <BlockTitle large>Favorites</BlockTitle>
    <List strongIos insetIos>
      {#each Object.entries(watchList) as [key, values] (key)}
        <ListItem link href={`/stop/${key}`} title={values} />
      {/each}
    </List>
  {:else}
    <Block class="flex flex-col items-center  gap-8 w-full pt-8">
      <div class="w-12 h-12 text-gray-300">
        <FaRegSmile />
      </div>
      <div class="flex items-center gap-2">
        <div class="w-4 h-4 text-yellow-500">
          <FaRegStar />
        </div>
        <p>Stations to add them to the watchlist</p>
      </div>
    </Block>
  {/if}
</div>
