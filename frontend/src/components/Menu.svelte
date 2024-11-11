<script>
    import AddSymbol from "./menu/AddSymbol.svelte";
    import BaseSize from "./menu/BaseSize.svelte";
    import PageSelectMenu from "./menu/PageSelectMenu.svelte";
    import { createEventDispatcher } from "svelte";

    let pageSelect;

    const dispatch = createEventDispatcher();

    const menuPageSelected = (e) => {
        let selectedPageNameString = e.detail;
        dispatch("menu-selected-page", selectedPageNameString);
    };
</script>

<main>
    <div class="container">
        <AddSymbol on:add-symbol />
        <BaseSize />
        <PageSelectMenu {pageSelect} on:page-select={menuPageSelected} />
    </div>
</main>

<style>
    .container {
        display: flex;
        gap: 2px;
        margin: 1rem 1rem 1rem 1rem;
        justify-content: right;
    }
</style>
