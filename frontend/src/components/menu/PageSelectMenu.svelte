<script>
    import { createEventDispatcher } from "svelte";
    let pageSelect;

    const dispatch = createEventDispatcher();

    const pageSelector = (pageString) => {
        dispatch("page-select", pageString);
    };
</script>

<main>
    <div class="container">
        <div class="btn_page" on:click={() => pageSelector("cards")}>cards</div>
        <div class="btn_page" on:click={() => pageSelector("activeTrades")}>trades</div>
    </div>
</main>

<style>
    .container {
        display: flex;
        padding: 12px;
        font-size: 14px;
        background: #1e222d;
    }

    .btn_page {
        user-select: none;

        display: flex;
        align-items: center;
        font-size: 14px;

        margin-right: 8px;
        padding-left: 6px;
        padding-right: 6px;

        height: 28px;
        text-align: center;

        background: #131722;
        border: 1px solid #464c5e;
        border-radius: 4px;
        color: #fefefe;

        cursor: pointer;
    }

    .btn_page:hover {
        background: #343d57;
    }
</style>
