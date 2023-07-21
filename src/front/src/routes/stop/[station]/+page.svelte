<script>
  import { Badge, Block, BlockHeader, BlockTitle, Card } from "konsta/svelte";
  import FaBus from "svelte-icons/fa/FaBus.svelte";
  import Departure from "../../../components/Departure.svelte";
</script>

<div class="flex flex-col">
  <div>
    <BlockTitle large>Viktoriaallee, Aachen</BlockTitle>
    <BlockHeader>
      <div class="flex flex-col gap-2">
        <div class="flex gap-2">
          <div class="w-4 h-4">
            <FaBus />
          </div>
          <Badge>33</Badge>
          <Badge>7</Badge>
        </div>
        <p class="text-xs">Last Update: 12:12</p>
      </div>
    </BlockHeader>
    <Block />
  </div>

  <div>
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
    <Departure />
  </div>
</div>

