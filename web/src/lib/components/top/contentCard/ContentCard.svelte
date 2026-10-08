<script lang="ts">
    import ContentCardImage from "#lib/components/top/contentCard/ContentCardImage.svelte";
    import ContentCardWrapper from "#lib/components/top/contentCard/ContentCardWrapper.svelte";
    import {cn} from "#lib/utils.ts";

    type Props = {
        title: string;
        contentID: string | null;
        pictureUrl: string | null;
        contentType: string;
        class?: string | null
    };

    let {title, contentID = $bindable(null), pictureUrl, contentType, class: className}: Props = $props();
</script>

<article class={cn('flex flex-col items-center justify-center gap-2 min-h-content min-w-24', className)}>
    <ContentCardWrapper {contentID} {contentType}>
        <ContentCardImage {contentType} {pictureUrl}/>
        <p class="mb-auto line-clamp-2 h-14 text-center md:text-lg">{title}</p>
    </ContentCardWrapper>
</article>
