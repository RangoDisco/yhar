<script lang="ts">
    import * as Command from "#lib/components/ui/command/index.ts";
    import * as Tooltip from "#lib/components/ui/tooltip/index.ts";
    import {Button} from "#lib/components/ui/button/index.ts";
    import {Spinner} from "#lib/components/ui/spinner/index.ts";
    import type {SearchResult} from "#lib/types/content.ts";
    import {debounce} from "#lib/utils.ts";
    import {page} from "$app/state";

    import Search from "@lucide/svelte/icons/search";
    import ContentPicture from "#lib/components/ContentPicture.svelte";
    import type {Paginated} from "#lib/types/pagination.ts";

    let open = $state(false);
    let searchQuery = $state("");
    let results: SearchResult[] | null = $state(null);
    let loading = $state(false);
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

    const fetchPage = async (query: string, page: number) => {
        loading = true;

        const res = await fetch(`/api/contents/search?value=${query}&page=${page}`);
        if (!res.ok) {
            hasError = true;
            loading = false;
            return;
        }

        const json: Paginated<SearchResult> = await res.json();
        if (!results) {
            results = json.results;
        } else {
            results = [...results, ...json.results];
        }

        hasNextPage = json.pagination.has_next_page;
        loading = false;
    };

    const handleSearch = debounce(async (query: string) => {
        results = null;
        hasError = false;
        currentPage = 1;

        if (query === "") {
            return;
        }

        await fetchPage(query, currentPage);
    }, 300);

    $effect(() => {
        if (!open) {
            searchQuery = "";
            results = null;
            hasNextPage = false;
            currentPage = 1;
            bottomListRef = null;
        }
    });

    $effect(() => {
        if (!bottomListRef || !hasNextPage || loading) {
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
    <Tooltip.Root>
        <Tooltip.Trigger>
            <Button onclick={() => (open = true)} variant="ghost" class="w-full" size="icon">
                <Search/>
            </Button>
        </Tooltip.Trigger>
        <Tooltip.Content side="right" align="center">
            <p>Search ⌘K</p>
        </Tooltip.Content>
    </Tooltip.Root>


    <Command.Dialog bind:open={open} shouldFilter={false}>
        <Command.Input
                placeholder="Search..."
                bind:value={searchQuery}
                oninput={(e) => handleSearch(e.currentTarget.value)}
        />
        <Command.List>
            {#if loading}
                <div class="flex items-center justify-center p-4">
                    <Spinner/>
                </div>
            {/if}

            {#if results && results.length > 0}
                <Command.Group heading="Results">
                    {#each results as result}
                        <Command.LinkItem
                                href={`/users/${page.params.userID}/top/${result.contentType}/${result.contentID}`}
                                class="cursor-pointer"
                                value={result.contentID}
                                onclick={() => (open = false)}
                        >
                            <ContentPicture
                                    class="h-8 w-8"
                                    url={`${result.domain}/images/${result.path}`}
                                    contentType={result.contentType}
                                    alt={result.title}
                            />
                            <span>{result.title}</span>
                        </Command.LinkItem>
                    {/each}
                    <div bind:this={bottomListRef}></div>
                </Command.Group>

            {:else if hasError}
                <Command.Empty color="red">An error occurred.</Command.Empty>
            {:else if results && results.length === 0}
                <Command.Empty>No results found.</Command.Empty>
            {/if}
        </Command.List>
    </Command.Dialog>
</div>
