<script>
  import { Badge, Block, BlockHeader, BlockTitle, Card } from "konsta/svelte";
  import FaBus from "svelte-icons/fa/FaBus.svelte";
  import Departure from "../../../components/Departure.svelte";

  import { page } from "$app/stores";
  import { onMount } from "svelte";

  let pageData = null;
  let now = null;

  onMount(async () => {
    let stopId = $page.params.station;
    const response = await fetch(`/api/departure/${stopId}`);
    const data = await response.json();
    pageData = data;
    now = new Date();
    now = now.toLocaleTimeString();
    console.log(pageData);
  });
</script>

<div class="flex flex-col">
  {#if pageData}
    <div>
      <BlockTitle large>{pageData.stop.name}</BlockTitle>
      <BlockHeader>
        <div class="flex flex-col gap-2">
          <div class="flex gap-2">
            <div class="w-4 h-4">
              <FaBus />
            </div>
            <Badge>33</Badge>
            <Badge>7</Badge>
          </div>
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
    <div>
      <BlockTitle large>Loading</BlockTitle>
    </div>
  {/if}
</div>
