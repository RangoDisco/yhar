<script lang="ts">
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import ContentCard from "$lib/components/top/contentCard/ContentCard.svelte";
    import UploadableImage from "$lib/components/top/UploadableImage.svelte";
    import HistoryList from "$lib/components/top/tracks/HistoryList.svelte";

    let {data} = $props();
</script>

<main class="flex flex-col gap-8">
    <section class="flex items-center gap-4">
        <UploadableImage pictureUrl={data.artist.picture_url} alt="{data.artist.name}'s picture" contentType="artists"
                         contentID={parseInt(data.artist.id)} uploadEnabled={data.user?.role === 'ADMIN'}/>
        <div class="flex flex-col">
            <p class="text-3xl">{data.artist.name}</p>
            <p class="text-muted-foreground">{data.artist.scrobble_count} scrobbles</p>
        </div>
    </section>

    <div class="flex flex-col gap-8">
        <!--TOP ALBUMS-->
        <ContentListWrapper title="Top albums" url="top/albums?artist={data.artist.id}">
            <div class="flex flex-wrap gap-4">
                {#each data.albums.results as album, i (album.id)}
                    <ContentCard title={album.title ?? album.name} pictureUrl={album.picture_url} contentType="albums"
                                 contentID={album.id}/>
                {/each}
            </div>
        </ContentListWrapper>

        <!--TOP TRACKS-->
        <ContentListWrapper title="Top tracks" url="top/tracks?artist={data.artist.id}">
            <div class="flex flex-col gap-4">
                {#each data.tracks.results as track, i (track.id)}
                    <ContentListItem index={i+1} title={track.title}
                                     pictureUrl={track.picture_url}
                                     scrobbleCount={track.scrobble_count}
                                     parentType="albums"
                                     parents={[track.album]}
                                     contentType="tracks"/>
                {/each}
            </div>
        </ContentListWrapper>

        <!--HISTORY-->
        <HistoryList scrobbles={data.history.results} url="history?artist={data.artist.id}" parentType="albums"/>
    </div>
</main>
