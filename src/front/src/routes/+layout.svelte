<script>
  import "../app.css";
  import { App, Page, Navbar } from "konsta/svelte";
</script>

<App theme="ios">
  <Page>
    <Navbar title="🚌 AC Bus Departures" />
    <slot />
  </Page>
</App>
