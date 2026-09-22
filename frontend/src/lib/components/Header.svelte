<script lang="ts">
  import { onMount } from "svelte";
  import { Window } from "@wailsio/runtime";
  import { Button } from "$lib/components/ui/button";
  import { Minus, Maximize, Minimize, X } from "@lucide/svelte";

  let maximized = false;

  onMount(async () => {
    maximized = await Window.IsMaximised();
  });

  function toggleMaximise() {
    Window.IsMaximised().then((isMaximised) => {
      if (isMaximised) {
        Window.UnMaximise();
      } else {
        Window.Maximise();
      }
      maximized = !isMaximised;
    });
  }
</script>

<header
  class="col-span-2 h-8 wails-drag border-b border-border flex items-center justify-between pl-4 bg-secondary"
>
  <div class="flex items-center gap-2 text-sm font-bold">Vlna</div>
  <div class="wails-no-drag flex">
    <Button variant="ghost" size="icon" onclick={Window.Minimise}>
      <Minus />
    </Button>
    <Button variant="ghost" size="icon" onclick={toggleMaximise}>
      {#if maximized}
        <Minimize />
      {:else}
        <Maximize />
      {/if}
    </Button>
    <Button variant="ghost" size="icon" onclick={Window.Close}>
      <X />
    </Button>
  </div>
</header>
