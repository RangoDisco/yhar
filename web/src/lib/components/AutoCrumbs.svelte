<script lang="ts">
    import {page} from "$app/state";
    import {
        Breadcrumb, BreadcrumbList, BreadcrumbItem,
        BreadcrumbLink, BreadcrumbSeparator, BreadcrumbPage
    } from "#lib/components/ui/breadcrumb/index.ts";

    type BreadcrumbItem = {
        href: string | null;
        label: string;
    };
    const base = `/users/${page.data.user?.id}`;

    const items = $derived.by(() => {
        const parts = page.url.pathname.split("/").filter(Boolean);
        const result: BreadcrumbItem[] = [{href: "/", label: "Stats"}, {
            href: `${base}/top`,
            label: page.data.user?.username
        }];

        if (parts.length > 3) {
            const part = parts[3];
            result.push({
                href: `${base}/top/${part}`,
                label: `${part.charAt(0).toUpperCase()}${part.slice(1)}`,
            });

            if (parts.length > 4) {
                switch (part) {
                    case "albums":
                        result.push({
                            href: `${base}/top/albums/${parts[4]}`,
                            label: page.data.album.title
                        });
                        break;
                    case "artists":
                        result.push({
                            href: `${base}/top/artists/${parts[4]}`,
                            label: page.data.artist.name,
                        });
                        break;
                    case "tracks":
                        result.push({
                            href: `${base}/top/tracks/${parts[4]}`,
                            label: page.data.track.title,
                        });
                        break;
                }
            }
        }

        return result;
    });
</script>

<Breadcrumb>
    <BreadcrumbList>
        {#each items as item, i}
            <BreadcrumbItem>
                {#if i < items.length - 1}
                    <BreadcrumbLink href={item.href}>{item.label}</BreadcrumbLink>
                    <BreadcrumbSeparator/>
                {:else}
                    <BreadcrumbPage>{item.label}</BreadcrumbPage>
                {/if}
            </BreadcrumbItem>
        {/each}
    </BreadcrumbList>
</Breadcrumb>
