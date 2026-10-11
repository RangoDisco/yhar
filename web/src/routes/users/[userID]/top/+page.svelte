<script lang="ts">
    import * as Tabs from "#lib/components/ui/tabs/index.js";
    import ContentListItem from "#lib/components/top/ContentListItem.svelte";
    import ContentListWrapper from "#lib/components/top/ContentListWrapper.svelte";
    import ContentCard from "#lib/components/top/contentCard/ContentCard.svelte";
    import {Period} from "#lib/types/period.js";
    import HistoryList from "#lib/components/top/tracks/HistoryList.svelte";
    import SvelteHead from "#lib/components/SvelteHead.svelte";
    import Heatmap from "#lib/components/heatmap/Heatmap.svelte";
    import StatCard from "#lib/components/top/StatCard.svelte";
    import dayjs from "dayjs";

    const {data} = $props();
    const periods = [Period.week, Period.month, Period.year, Period.overall];
</script>

<SvelteHead title="Top"/>
<div class="relative grid grid-cols-4 gap-12">
    <div class="flex flex-col gap-8 col-span-4 2xl:col-span-3">
        <h1 class="text-2xl text-muted-foreground">Recap</h1>
        <!--All time recap-->
        <div class="flex flex-col gap-8 lg:flex-row lg:justify-between">
            <div class="grid grid-cols-2 gap-6 lg:w-1/3">
                <StatCard title="Total listens" description={data.totalScrobblesStats.totalScrobbles.toString()}/>
                <StatCard title="Total time"
                          description={dayjs(data.totalScrobblesStats.totalDuration).format('DD [days]')}/>
                <StatCard title="Avg listens per day" description={data.avgDailyScrobbles.avgCount.toFixed(2)}/>
                <StatCard title="Tracks listened" description={data.totalDistinctTracksCount.totalCount.toString()}/>
            </div>
            <div class="flex flex-col overflow-scroll lg:w-2/3">
                <span class="font-medium text-muted-foreground mb-2">{data.yearlyScrobbleCount}
                    scrobbles in the last year</span>
                <Heatmap data={data.timeline}/>
            </div>
        </div>
        <Tabs.Root value="week">
            <Tabs.List class="w-full mx-auto">
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
                        <ContentListWrapper title="Artists" url="top/artists">
                            <div class="grid gap-4 grid-cols-3 md:grid-flow-col md:grid-cols-[repeat(auto-fill,minmax(160px,1fr))] md:auto-cols-[minmax(160px,1fr)] overflow-auto">
                                {#each periodData?.artists.results as artist, i (artist.id)}
                                    <div class={i >= 6 ? 'hidden lg:block' : ''}>
                                        <ContentCard
                                                title={artist.name}
                                                pictureUrl={artist.picture_url}
                                                contentType="artists"
                                                contentID={artist.id}
                                        />
                                    </div>
                                {/each}
                            </div>
                        </ContentListWrapper>

                        <!--TOP ALBUMS-->
                        <ContentListWrapper title="Albums" url="top/albums">
                            <div class="grid gap-4 grid-cols-3 md:grid-flow-col md:grid-cols-[repeat(auto-fill,minmax(160px,1fr))] md:auto-cols-[minmax(160px,1fr)] overflow-auto">
                                {#each periodData?.albums.results as album, i (album.id)}
                                    <div class={i >= 6 ? 'hidden lg:block' : ''}>
                                        <ContentCard
                                                title={album.title}
                                                pictureUrl={album.picture_url}
                                                contentType="albums"
                                                contentID={album.id}
                                        />
                                    </div>
                                {/each}
                            </div>
                        </ContentListWrapper>

                        <!--TOP TRACKS-->
                        <ContentListWrapper title="Tracks" url="top/tracks">
                            <div class="flex flex-col gap-4">
                                {#each periodData?.tracks.results as track, i (track.id)}
                                    <ContentListItem
                                            index={i + 1}
                                            title={track.title}
                                            pictureUrl={track.picture_url}
                                            scrobbleCount={track.scrobble_count}
                                            parentType="artists"
                                            parents={track.artists}
                                            contentType="tracks"
                                    />
                                {/each}
                            </div>
                        </ContentListWrapper>
                    {/await}
                </Tabs.Content>
            {/each}
        </Tabs.Root>
    </div>

    <!--HISTORY-->
    <div class="sticky col-span-4 top-5 h-fit md:-block 2xl:col-span-1">
        <HistoryList scrobbles={data.history.results} url="history" parentType="artists"
                     shouldShowDeleteButton={false}/>
    </div>
</div>
