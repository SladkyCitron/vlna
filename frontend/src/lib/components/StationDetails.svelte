<script lang="ts">
  import * as Sheet from "$lib/components/ui/sheet";
  import {
    selectedStation,
    isStationDetailsOpen,
    closeStationDetails,
  } from "$lib/stores/stationDetails";
  import * as m from "$lib/paraglide/messages.js";

  function handleClose(open: boolean) {
    if (!open) {
      open = false;
      setTimeout(() => {
        closeStationDetails();
      }, 200);
    }
  }
</script>

<Sheet.Root bind:open={$isStationDetailsOpen} onOpenChange={handleClose}>
  <Sheet.Content side="right">
    <Sheet.Header>
      <Sheet.Title>{m.stationDetails()}</Sheet.Title>
    </Sheet.Header>
  </Sheet.Content>
</Sheet.Root>
