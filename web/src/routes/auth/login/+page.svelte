<script lang="ts">
    import GalleryVerticalEndIcon from "@lucide/svelte/icons/gallery-vertical-end";
    import type {HTMLAttributes} from "svelte/elements";
    import {
        FieldGroup,
        Field,
        FieldLabel, FieldError,
    } from "$lib/components/ui/field/index.js";
    import {Input} from "$lib/components/ui/input/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import type {PageProps} from "./$types";


    let {data, form}: PageProps = $props();
    console.log(form);
</script>

<div class="w-md md:w-lg 2xl:w-xl max-w-2xl mx-auto justify-center">
    <form method="POST">
        <FieldGroup>
            <div class="flex flex-col items-center gap-2 text-center">
                <a href="##" class="flex flex-col items-center gap-2 font-medium">
                    <div class="flex size-8 items-center justify-center rounded-md">
                        <GalleryVerticalEndIcon class="size-6"/>
                    </div>
                    <span class="sr-only">Yhar</span>
                </a>
                <h1 class="text-xl font-bold">Yhar</h1>
            </div>
            <Field>
                <FieldLabel for="username">Username</FieldLabel>
                <Input name="username" id="username" type="text" required value={form?.username ?? ''}/>
                {#if form?.field === "username"}
                    <FieldError>{form?.error}</FieldError>
                {/if}
            </Field>
            <Field>
                <FieldLabel for="password">Password</FieldLabel>
                <Input name="password" id="password" type="password" required value={form?.password ?? ''}/>
                {#if form?.field === "password"}
                    <FieldError>{form?.error}</FieldError>
                {/if}
            </Field>
            <Field>
                <Button type="submit" class="cursor-pointer">Login</Button>
                {#if form?.field === null && form?.error}
                    <FieldError>{form?.error}</FieldError>
                {/if}
            </Field>
        </FieldGroup>
    </form>
</div>
