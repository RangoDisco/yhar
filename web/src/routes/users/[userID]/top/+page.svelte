<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import ContentCard from "$lib/components/top/contentCard/ContentCard.svelte";
    import {Period} from "$lib/types/period";
    import HistoryList from "$lib/components/top/tracks/HistoryList.svelte";

    const {data} = $props();
    const periods = [Period.week, Period.month, Period.year, Period.overall];
</script>
<div class="flex flex-col gap-8">
    <Tabs.Root value="week">
        <Tabs.List class="w-full">
            {#each periods as period}
                <Tabs.Trigger value={period}>{period.charAt(0).toUpperCase() + period.slice(1)}</Tabs.Trigger>
            {/each}
        </Tabs.List>
        {#each periods as period}
            <Tabs.Content value={period} class="flex flex-col gap-8">
                {#await data[period]}
                    Loading...
                {:then periodData}
                    <!--TOP ARTISTS-->
                    <ContentListWrapper title="Top artists" url="top/artists">
                        <div class="flex flex-wrap gap-4">
                            {#each periodData?.artists.results as artist, i (artist.id)}
                                <ContentCard title={artist.name} pictureUrl={artist.picture_url} contentType="artists"
                                             contentID={artist.id}/>
                            {/each}
                        </div>
                    </ContentListWrapper>

                    <!--TOP ALBUMS-->
                    <ContentListWrapper title="Top albums" url="top/albums">
                        <div class="flex flex-wrap gap-4">
                            {#each periodData?.albums.results as album, i (album.id)}
                                <ContentCard title={album.title} pictureUrl={album.picture_url}
                                             contentID={album.id}
                                             contentType="albums"/>
                            {/each}
                        </div>
                    </ContentListWrapper>

                    <!--TOP TRACKS-->
                    <ContentListWrapper title="Top tracks" url="top/tracks">
                        <div class="flex flex-col gap-4">
                            {#each periodData?.tracks.results as track, i (track.id)}
                                <ContentListItem index={i+1} title={track.title}
                                                 pictureUrl={track.picture_url}
                                                 scrobbleCount={track.scrobble_count}
                                                 parentType="artists"
                                                 parents={track.artists}
                                                 contentType="tracks"/>
                            {/each}
                        </div>
                    </ContentListWrapper>
                {/await}
            </Tabs.Content>
        {/each}
    </Tabs.Root>

    <!--HISTORY-->
    <HistoryList scrobbles={data.history.results} url="history" parentType="artists"/>
</div>
