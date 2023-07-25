<script>
  import { Card } from "konsta/svelte";
  export let departure;
  let timePlanned;
  let timePrognosed;

  $: {
    // result comes in format 'HHMMSS'
    // convert to 'HH:MM'

    timePlanned =
      departure.planned !== ""
        ? departure.planned.slice(0, 2) + ":" + departure.planned.slice(2, 4)
        : "";
    timePrognosed =
      departure.prognosed !== ""
        ? departure.prognosed.slice(0, 2) +
          ":" +
          departure.prognosed.slice(2, 4)
        : "";
  }
</script>

<Card>
  <div class="flex gap-4 items-center">
    <div class="flex flex-col">
      <p class={`text-xl ${departure.cancelled && "line-through"}`}>
        {timePlanned}
      </p>
      <p class={`text-red-700 ${departure.cancelled && "line-through"}`}>
        {timePrognosed}
      </p>
    </div>
    <div class="flex flex-col grow">
      <p class={`text-xl ${departure.cancelled && "line-through"}`}>
        {departure.bus}
      </p>
      <p>
        {#if departure.cancelled}
          cancelled
        {:else}
          &#8594; {departure.destination}
        {/if}
      </p>
    </div>
    <div class={`text-2xl ${departure.cancelled && "line-through"}`}>
      {departure.platform}
    </div>
  </div>
</Card>
