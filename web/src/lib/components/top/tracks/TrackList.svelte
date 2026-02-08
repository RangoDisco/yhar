<script>
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";

    let {
        tracks
    } = $props();
</script>


<ContentListWrapper title="Top tracks" url="/tracks">
    <div class="flex flex-col gap-4">
        {#each tracks as track, i (track.id)}
            <ContentListItem index={i} title={track.title} parents={track.artists}
                             pictureUrl={track.picture_url}
                             scrobbleCount={track.scrobble_count}
                             contentType="tracks"/>
        {/each}
    </div>
</ContentListWrapper>
