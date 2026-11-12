<script lang="ts">
    import * as Command from "#lib/components/ui/command/index.ts";
    import * as Kbd from "#lib/components/ui/kbd/index.js";
    import * as Sidebar from "#lib/components/ui/sidebar/index.ts";
    import * as Tabs from "#lib/components/ui/tabs/index.ts";
    import {Spinner} from "#lib/components/ui/spinner/index.ts";
    import type {SearchResult} from "#lib/types/content.ts";
    import {debounce} from "#lib/utils.ts";
    import {page} from "$app/state";

    import Search from "@lucide/svelte/icons/search";
    import ContentPicture from "#lib/components/ContentPicture.svelte";
    import type {Paginated} from "#lib/types/pagination.ts";

    let open = $state(false);
    let searchQuery = $state("");
    let contentType: string | null = $state(null);
    let results: SearchResult[] | null = $state(null);
    let isLoading = $state(false);
    let isLoadingMore = $state(false);
    let hasError = $state(false);
    let hasNextPage = $state(false);
    let currentPage = $state(1);
    let bottomListRef: HTMLDivElement | null = $state(null);

    const handleKeyDown = (event: KeyboardEvent) => {
        if (event.key === "k" && (event.metaKey || event.ctrlKey)) {
            event.preventDefault();
            open = true;
        }
    };

    const fetchPage = async (query: string, page: number, type?: string | null = null) => {
        isLoadingMore = true;
        contentType = null;

        let url = `/api/contents/search?value=${query}&page=${page}`;
        if (type) {
            url += `&type=${type}`;
        }

        const res = await fetch(url);
        if (!res.ok) {
            hasError = true;
            isLoadingMore = false;
            return;
        }

        const json: Paginated<SearchResult> = await res.json();
        if (!results) {
            results = json.results;
        } else {
            results = [...results, ...json.results];
        }

        contentType = type;
        hasNextPage = json.pagination.has_next_page;
        isLoadingMore = false;
    };

    const handleSearch = debounce(async (query: string, type?: string | null) => {
        isLoading = true;
        results = null;
        hasError = false;
        currentPage = 1;

        if (query === "") {
            isLoading = false;
            return;
        }

        await fetchPage(query, currentPage, type);
        isLoading = false;
    }, 300);

    // EFFECTS
    $effect(() => {
        if (!open) {
            searchQuery = "";
            contentType = null;
            results = null;
            hasNextPage = false;
            currentPage = 1;
            bottomListRef = null;
            isLoading = false;
            isLoadingMore = false;
        }
    });

    $effect(() => {
        if (!bottomListRef || !hasNextPage || isLoading || isLoadingMore) {
            return;
        }

        const observer = new IntersectionObserver(async (entries) => {
            if (entries[0].isIntersecting) {
                currentPage += 1;
                await fetchPage(searchQuery, currentPage);
            }
        });

        observer.observe(bottomListRef);

        return () => {
            observer.disconnect();
        };
    });
</script>

<svelte:document onkeydown={handleKeyDown}/>

<div class="flex flex-col gap-4">
    <Sidebar.MenuButton
            onclick={() => (open = true)}>
        {#snippet tooltipContent()}
            <div>
                Search
                <Kbd.Root>⌘K</Kbd.Root>
            </div>
        {/snippet}
        <Search/>
        <span>Search</span>
    </Sidebar.MenuButton>

    <Command.Dialog bind:open={open} shouldFilter={false} class="top-1/4 sm:max-w-md">
        <Command.Input
                placeholder="Search..."
                bind:value={searchQuery}
                oninput={(e) => handleSearch(e.currentTarget.value)}
        />
        <Command.List class="max-h-100">
            {#if (results && results.length > 0) || isLoading}
                <Tabs.Root value={contentType} class="w-full mt-3 px-2" activationMode="manual"
                           onValueChange={(value) => handleSearch(searchQuery, value)}>
                    <Tabs.List>
                        <Tabs.Trigger value={null}>All</Tabs.Trigger>
                        <Tabs.Trigger value="artists">Artists</Tabs.Trigger>
                        <Tabs.Trigger value="albums">Albums</Tabs.Trigger>
                        <Tabs.Trigger value="tracks">Tracks</Tabs.Trigger>
                    </Tabs.List>
                </Tabs.Root>
            {/if}

            {#if isLoading}
                <div class="flex items-center justify-center p-4">
                    <Spinner/>
                </div>
            {/if}

            {#if !isLoading && results && results.length > 0}
                <Command.Group heading="Results">
                    {#each results as result}
                        <Command.LinkItem
                                href={`/users/${page.params.userID}/top/${result.contentType}/${result.contentID}`}
                                class="cursor-pointer max-h-14"
                                value={result.contentID}
                                onclick={() => (open = false)}
                        >
                            <ContentPicture
                                    class="h-8 min-w-8"
                                    url={`${result.domain}/images/${result.path}`}
                                    contentType={result.contentType}
                                    alt={result.title}
                            />
                            <div class="flex flex-col gap-1/2">
                                <span class="text-base line-clamp-1">{result.title}</span>
                                <span class="text-sm font-medium text-muted-foreground">{result.contentType.charAt(0).toUpperCase() + result.contentType.slice(1)}</span>
                            </div>
                        </Command.LinkItem>
                    {/each}
                    <div bind:this={bottomListRef}></div>

                    {#if isLoadingMore}
                        <div class="flex items-center justify-center p-4">
                            <Spinner/>
                        </div>
                    {/if}
                </Command.Group>

            {:else if hasError}
                <Command.Empty color="red">An error occurred.</Command.Empty>
            {:else if results && results.length === 0}
                <Command.Empty>No results found.</Command.Empty>
            {/if}
        </Command.List>
    </Command.Dialog>
</div>
