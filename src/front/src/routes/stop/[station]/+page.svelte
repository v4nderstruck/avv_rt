<script>
  import { Badge, Block, BlockHeader, BlockTitle, Card } from "konsta/svelte";
  import FaBus from "svelte-icons/fa/FaBus.svelte";
  import Departure from "../../../components/Departure.svelte";
  import FaRegStar from "svelte-icons/fa/FaRegStar.svelte";
  import { page } from "$app/stores";
  import { onMount } from "svelte";

  let pageData = null;
  let stopId;
  let now = null;
  let localStorage = null;

  let watchList = {};

  onMount(async () => {
    localStorage = window.localStorage;
    if (localStorage && localStorage.getItem("stationWatchList")) {
      watchList = JSON.parse(localStorage.getItem("stationWatchList"));
    }
    stopId = $page.params.station;
    const response = await fetch(`/api/departure/${stopId}/30`);
    const data = await response.json();
    pageData = data;
    now = new Date();
    now = now.toLocaleTimeString();
  });

  $: isOnWatchList = () => {
    return stopId in watchList;
  };

  const favClickHandler = () => {
    if (isOnWatchList()) {
      const { [stopId]: something, ...rest } = watchList;
      watchList = rest;
      console.log("remove ", stopId, watchList);
    } else {
      watchList = { ...watchList, [stopId]: pageData.stop.name };
      console.log("add ", watchList);
    }
    if (localStorage) {
      localStorage.setItem("stationWatchList", JSON.stringify(watchList));
    }
  };
</script>

<div class="flex flex-col">
  {#if pageData}
    <div>
      <BlockTitle large>
        <div class="flex w-full justify-between items-center">
          <p>{pageData.stop.name}</p>
          <button
            on:click={favClickHandler}
            class="w-8 h-8"
            class:text-yellow-400={isOnWatchList()}
          >
            <FaRegStar />
          </button>
        </div>
      </BlockTitle>
      <BlockHeader>
        <div class="flex flex-col gap-2">
          <!-- <div class="flex gap-2"> -->
          <!--   <div class="w-4 h-4"> -->
          <!--     <FaBus /> -->
          <!--   </div> -->
          <!--   <Badge>33</Badge> -->
          <!--   <Badge>7</Badge> -->
          <!-- </div> -->
          <p class="text-xs">Last Update: {now}</p>
        </div>
      </BlockHeader>
      <Block />
    </div>
    <div>
      {#each pageData.departures as dep (dep.journey + dep.planned)}
        <Departure departure={dep} />
      {/each}
    </div>
  {:else}
    <div class="w-screen h-screen flex items-center justify-center">
      <BlockTitle large>Loading...</BlockTitle>
    </div>
  {/if}
</div>
