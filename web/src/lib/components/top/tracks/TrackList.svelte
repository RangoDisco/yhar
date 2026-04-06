<script lang="ts">
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import type {Track} from "$lib/types/content";

    type Props = {
        tracks: Track[]
        params: string | null
    }

    let {
        tracks,
        params = $bindable(""),
    }: Props = $props();
</script>


<ContentListWrapper title="Top tracks" url="top/tracks{params}">
    <div class="flex flex-col gap-4">
        {#each tracks as track, i (track.id)}
            <ContentListItem index={i} title={track.title}
                             pictureUrl={track.picture_url}
                             scrobbleCount={track.scrobble_count}
                             parentType="artists"
                             parents={track.artists}
                             contentType="tracks"/>
        {/each}
    </div>
</ContentListWrapper>
