<script lang="ts">
    import ContentCardImage from "$lib/components/top/contentCard/ContentCardImage.svelte";
    import ContentCardWrapper from "$lib/components/top/contentCard/ContentCardWrapper.svelte";

    let {
        title,
        contentID = $bindable(null),
        pictureUrl,
        contentType
    } = $props();

</script>

<article class="w-[30%] md:w-[15%] flex flex-col items-center justify-center gap-2">
    <ContentCardWrapper {contentID} {contentType}>
        <ContentCardImage {contentType} {contentID} {pictureUrl}/>
        <p class="text-center mb-auto">{title}</p>
    </ContentCardWrapper>
</article>
