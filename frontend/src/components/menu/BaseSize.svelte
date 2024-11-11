<script>
    import { onMount } from "svelte";
    import { baseUrl } from "../baseUrl";

    let baseSize = 0;
    let error;

    onMount(async () => {
        try {
            const res = await fetch(`${baseUrl}/baseSize/get`);
            let value = await res.json();
            baseSize = value.BaseSize;
        } catch (e) {
            error = e;
        }
    });

    const setBaseSize = async () => {
        const msgBody = {
            BaseSize: Math.trunc(Number(baseSize)),
        };

        const url = `${baseUrl}/baseSize/set`;

        await fetch(url, {
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
            method: "POST",
            body: JSON.stringify(msgBody),
        });
    };
</script>

<main>
    <div class="container">
        <div class="text">Base Size</div>
        <form on:submit|preventDefault={setBaseSize}>
            <div class="base-size-input">
                <input type="text" bind:value={baseSize} />
                <i>$</i>
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

        padding-left: 14px;
        max-width: 48px;
        height: 28px;

        background: #131722;
        border: 1px solid #464c5e;
        border-radius: 4px;
        color: #fefefe;
    }

    .base-size-input {
        position: relative;
    }

    .base-size-input > i {
        position: absolute;
        display: block;
        transform: translate(0, -50%);
        top: 50%;
        pointer-events: none;
        width: 20px;
        text-align: center;
        font-style: normal;
        font-size: 14px;
        color: #cacaca;
    }
</style>
