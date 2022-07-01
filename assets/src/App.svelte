<script>
    import "node_modules/json2/lib/JSON2/static/json2.js";
    import "./css/fontawesome.css";
    import Share from './Share.svelte';
    let name = 'DHT crawler';

    let init = true;
    let resultsPresent = false;
    let pn = 1;
    let limit = 10;
    let resultsFound = 0;
    let APILocation = window.location.host;
    let searchQuery = "";
    let searchResult = [];
    
    $: if (!init) { 
        displayPage(searchQuery, limit, pn);
    } else {
        init = false;
    }
    $: pages = genPageslist(pn, resultsFound, limit);
    function genPageslist(pn, resultsFound, limit) {
        let pagesQuantity = Math.ceil(resultsFound/limit);
        let res = [];
        res.length = 0;
        if (pn < 3) {
            res = [1,2,3,4,5];
        } else {
            res = [pn-2, pn-1, pn, pn+1, pn+2];
        }
        while(res[res.length-1] > pagesQuantity) {
            res.pop();
        } 
        return res;
    }
    function nextPage() {pn++;}
    function previousPage() {if (pn > 1) pn--;}
    function changePage(p) {
        pn = p;
    }
    function submitQuery(e) {
        if (e.target[0].value != "") {
            searchQuery = e.target[0].value;
            pn = 1;
        } else {
            searchQuery = e.target[0].value;
            resultsPresent = false;
        }
    }
    async function displayPage(query, q, pn) {
        if (query.trim() == "") 
            return;
        let reqString = `http://${APILocation}/api/v1/dhtcrawler/displaypage?size=${limit}&n=${pn}&query=${query}`
        let response = await fetch(reqString);
        console.log(`reqString : ${reqString}`);
        console.log(`code : ${response.status}`);
        if (response.ok) {
            let respObj = await response.json();
            resultsFound = parseInt(respObj.Total, 10);
            searchResult =  respObj.Results;
            if (searchResult != "") {
                resultsPresent = true;
            } else {
                alert("Nothing found.");
            }
        } else {
            alert("HTTP Err : " + response.status);
        }
    }
</script>

<style>
    @font-face {
        font-family: "Font Awesome 6 Free";
        src: url("/webfonts/fa-solid-900.woff2") format("woff2");
    }

    :global(body) {
        padding: 0;
        font-family: Verdana, Arial, Helvetica, sans-serif;
        overflow: auto;
    }       
            
    #main {
        background-color: #f5ef9b;
        border-left: solid;
        border-right: solid;

        display: block;
        margin: auto;

        padding: 1pt;
        min-height: 100%;
    }       
    
    @media screen and (min-width: 800px) {            
        #main {
            width: 60%;
        }
    }
    
    #main > * {
        margin: 1%; 
        display: block;
    }

    button,h1,fieldset,form,#nav {
        text-align: center;
        color: black;
    }

    fieldset>label,h5 {
        display: inline-block;
        margin-top: 0;
        margin-bottom: 0;
    }
    
    form>#searchbar {
        width: 50%; 
        height: 37px;
    }

    form>#searchbutton {
        height: 37px;
    }

    #nav>button {
        margin: 2pt;
    }

    .highlighted {
        background-color: #dbdbda;
    }
</style>

<svelte:head>
    <title>{name}</title>
</svelte:head>
<div id="main">
    <h1>{name}</h1>
    <form on:submit|preventDefault={submitQuery} autocomplete="on" accept-charset="utf-8">
        <input type="text" id="searchbar" value=""/>
        <button id="searchbutton"> 
            <i class="fa fa-solid fa-magnifying-glass"></i>
        </button>
    </form>
        
    <fieldset>
        <legend>Limit</legend>
        <label><input type="radio" name="radio" on:click={() => {limit = 10;}} checked> 10 </label>
        <label><input type="radio" name="radio" on:click={() => {limit = 20;}}> 20 </label>
        <label><input type="radio" name="radio" on:click={() => {limit = 50;}}> 50 </label>
    </fieldset>
    
    <div id="info">
        {#if resultsPresent}
            <h5>Page: {pn}</h5>
            <h5>-</h5>
            <h5>Found in total: {resultsFound}</h5>
        {/if}
    </div>
    
    <hr>
    <div id="results">
        {#if resultsPresent}
            {#each searchResult as d, i}
                <Share no={i+1} {...d} />
            {/each}
        {/if}
    </div>
    <hr>
    
    {#if resultsPresent}
        <div id="nav">            
            <button on:click={previousPage}>&lt;</button>
            {#each pages as page}
                {#if page != pn}
                    <button on:click={changePage(page)}>{page}</button>
                {:else}
                    <button on:click={changePage(page)} class="highlighted">{page}</button>
                {/if}
            {/each}
            {#if genPageslist(pn, resultsFound, limit).length > 3}
                <button on:click={nextPage}>&gt;</button>
            {/if}
        </div>
    {/if}
</div>
