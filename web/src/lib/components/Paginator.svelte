<script>
    import * as Pagination from "$lib/components/ui/pagination/index";
    import {handlePageChange, PER_PAGE} from "$lib/pagination";

    let {totalCount, page} = $props();

</script>

<Pagination.Root count={totalCount} perPage={PER_PAGE} page={page}
                 onPageChange={handlePageChange}>
    {#snippet children({pages, currentPage})}
        <Pagination.Content>
            <Pagination.Item>
                <Pagination.Previous/>
            </Pagination.Item>
            {#each pages as page (page.key)}
                {#if page.type === "ellipsis"}
                    <Pagination.Item>
                        <Pagination.Ellipsis/>
                    </Pagination.Item>
                {:else}
                    <Pagination.Item>
                        <Pagination.Link {page} isActive={currentPage === page.value}>
                            {page.value}
                        </Pagination.Link>
                    </Pagination.Item>
                {/if}
            {/each}
            <Pagination.Item>
                <Pagination.Next/>
            </Pagination.Item>
        </Pagination.Content>
    {/snippet}
</Pagination.Root>
