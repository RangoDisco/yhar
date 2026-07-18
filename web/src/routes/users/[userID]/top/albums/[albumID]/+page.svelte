<script lang="ts">
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import {page} from "$app/state";
    import UploadableImage from "$lib/components/top/UploadableImage.svelte";
    import SvelteHead from "$lib/components/SvelteHead.svelte";

    let {data} = $props();
</script>


<SvelteHead title={data.album.title} description="Top tracks for {data.album.title}"/>
<main class="flex flex-col gap-8">
    <section class="flex items-center gap-4">
        <UploadableImage pictureUrl={data.album.picture_url} alt="{data.album.title}'s picture" contentType="albums"
                         contentID={parseInt(data.album.id)}
                         uploadEnabled={data.user?.role === 'ADMIN'}/>
        <div class="flex flex-col">
            <p class="text-3xl lg:text-4xl">{data.album.title}</p>
            <div class="flex gap-1 items-center">
                {#each data.album.artists as artist, i}
                    {#if i < 3}
                        {#if i !== 0}
                            ·
                        {/if}
                        <a class="text-sm text-muted-foreground hover:underline"
                           href="/users/{page.params.userID}/top/artists/{artist.id}">{artist.name}</a>
                    {/if}
                {/each}
            </div>
            <p class="text-muted-foreground">{data.album.scrobble_count} scrobbles</p>
        </div>
    </section>

    <div class="flex flex-col gap-8">
        <!--TOP TRACKS-->
        <ContentListWrapper title="Top tracks">
            <div class="flex flex-col gap-4">
                {#each data.tracks.results as track, i (track.id)}
                    <ContentListItem index={i+1} title={track.title}
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
