<script lang="ts">
    import ContentCardImage from "$lib/components/top/contentCard/ContentCardImage.svelte";
    import ContentCardWrapper from "$lib/components/top/contentCard/ContentCardWrapper.svelte";
    import {cn} from "$lib/utils.ts";

    type Props = {
        title: string
        contentID: string | null
        pictureUrl: string | null
        contentType: string
    }

    let {
        title,
        contentID = $bindable(null),
        pictureUrl,
        contentType,
    }: Props = $props();

</script>

<article class={cn("flex flex-col items-center justify-center gap-2")}>
    <ContentCardWrapper {contentID} {contentType}>
        <ContentCardImage {contentType} {pictureUrl}/>
        <p class="text-center mb-auto h-12 line-clamp-2 md:text-lg">{title}</p>
    </ContentCardWrapper>
</article>
