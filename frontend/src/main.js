import "./app.css";
import App from "./App.svelte";

const app = new App({
    target: document.getElementById("app"),
});

export default app;

/*  
    ----- ID LIST -----
    
    Settings:
        Side=Long       set_Side_Long
        Side=Short      set_Side_Short

        input Entry     input_Entry_#
        input SL        input_SL_#
        input TP        input_TP_#

        EntryType=Usd   set_EntryType_Usd
        EntryType=Pct   set_EntryType_Pct

    Buttons:
        SetOrder        btn_SetOrder_#
        CancelAllOrders btn_CancelAllOrders_#
        MarketEntry     btn_MarketEntry_#
        MarketExit      btn_MarketExit_#

    https://www.webmound.com/create-keyboard-shortcuts-in-javascript/

*/
