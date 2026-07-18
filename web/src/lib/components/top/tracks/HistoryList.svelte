<script lang="ts">

    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
    import type {Scrobble} from "$lib/types/content";
    import {toast} from "svelte-sonner";

    type Props = {
        scrobbles: Scrobble[]
        url: string | null;
        parentType: "artists" | "albums"
    }

    let {scrobbles, url = $bindable(null), parentType}: Props = $props();

    const handleDelete = async (id: string) => {
        const initialState = scrobbles;
        scrobbles = scrobbles.filter((scrobble) => scrobble.id !== id);
        const res = await fetch(`/api/scrobbles/${id}`, {
            method: "DELETE"
        });
        if (!res.ok) {
            scrobbles = initialState;
            toast.error("Unable to delete scrobble");
            return;
        }
        toast.success("Scrobble deleted successfully!");
    };
</script>

<ContentListWrapper title="HISTORY" url={url}>
    <div class="flex flex-col gap-2 w-full">
        {#each scrobbles as scrobble}
            <HistoryListItem scrobble={scrobble} parentType={parentType} handleDelete={handleDelete}/>
        {/each}
    </div>
</ContentListWrapper>
