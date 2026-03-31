<script lang="ts">
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";

    let {data} = $props();
</script>

<main class="flex flex-col gap-8">
    <section class="flex items-center gap-4">
        <img class="rounded-md aspect-square h-24" src={data.album.picture_url} alt="{data.album.title}'s picture"/>
        <div class="flex flex-col">
            <p>{data.album.title}</p>
            <p class="text-muted-foreground">{data.album.scrobble_count} scrobbles</p>
        </div>
    </section>

    <div class="flex flex-col gap-8">
        <!--TOP TRACKS-->
        <ContentListWrapper title="Top tracks">
            <div class="flex flex-col gap-4">
                {#each data.tracks.results as track, i (track.id)}
                    <ContentListItem index={i} title={track.title}
                                     pictureUrl={track.picture_url}
                                     scrobbleCount={track.scrobble_count}
                                     parentType="albums"
                                     parents={[track.album]}
                                     contentType="tracks"
                    />
                {/each}
            </div>
        </ContentListWrapper>
    </div>
</main>
