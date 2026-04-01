<script lang="ts">
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import {page} from "$app/state";

    let {data} = $props();
</script>

<main class="flex flex-col gap-8">
    <section class="flex items-center gap-4">
        <img class="rounded-md aspect-square h-24" src={data.album.picture_url} alt="{data.album.title}'s picture"/>
        <div class="flex flex-col">
            <p class="text-3xl">{data.album.title}</p>
            {#each data.album.artists as artist, i}
                {#if i < 3}
                    {#if i !== 0}
                        ·
                    {/if}
                    <a class="text-sm text-muted-foreground hover:underline"
                       href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a>
                {/if}
            {/each}
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
                                     parentType="artists"
                                     parents={track.artists}
                                     contentType="tracks"
                    />
                {/each}
            </div>
        </ContentListWrapper>
    </div>
</main>
