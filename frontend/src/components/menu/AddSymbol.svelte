<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { baseUrl } from "../baseUrl";

    let symbol = "";
    let listOfFutures = [];

    const dispatch = createEventDispatcher();

    let validationPassed = true;

    onMount(async () => {
        try {
            const res = await fetch(`${baseUrl}/listFuturesSymbols`);
            listOfFutures = await res.json();
        } catch (e) {
            error = e;
        }
    });

    const addSymbolReq = async () => {
        symbol = symbol.toUpperCase();

        let symbolFound = false;

        for (const i in listOfFutures) {
            if (symbol == listOfFutures[i]) {
                symbolFound = true;
                break;
            }
        }

        if (symbolFound) {
            validationPassed = true;

            dispatch("add-symbol", symbol);
            symbol = "";
        } else {
            validationPassed = false;
        }
    };
</script>

<main>
    <div class="container">
        <div class="text">Add Symbol</div>
        <form on:submit|preventDefault={addSymbolReq}>
            <div class="base-size-input">
                {#if !validationPassed}
                    <div class="verification-hint">Wrong name</div>
                {/if}
                <input type="text" bind:value={symbol} />
            </div>
        </form>
    </div>
</main>

<style>
    .container {
        user-select: none;

        display: flex;
        padding: 12px;
        font-size: 14px;
        background: #1e222d;
    }

    .text {
        margin-top: 7px;
        margin-right: 16px;
    }

    input {
        font-size: 14px;

        padding-left: 6px;
        max-width: 80px;
        height: 28px;

        background: #131722;
        border: 1px solid #464c5e;
        border-radius: 4px;
        color: #fefefe;
    }

    .verification-hint {
        position: absolute;
        transform: translate(-5%, -120%);
        color: white;
        background-color: red;
        padding: 8px;
        border-radius: 8px;
        opacity: 0.8;
    }
</style>
