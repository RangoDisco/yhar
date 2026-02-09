<script lang="ts">
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";

    type Props = {
        tracks: [],
        mode: "artists" | "albums"
        params: string | null
    }

    let {
        tracks,
        mode,
        params = $bindable(""),
    }: Props = $props();
</script>

<ContentListWrapper title="History" url={`history${params}`}>
    <div class="flex flex-col gap-2">
        {#each tracks as track}
            <HistoryListItem track={track} mode={mode}/>
        {/each}
    </div>
</ContentListWrapper>
