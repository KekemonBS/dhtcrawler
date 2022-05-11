<script>
    let name = 'DTH crawler';
    let resultsPresent = true;
    let pn = 1;
    let limit = 10;
    let resultsFound = 666;
    //$: displayPage(limit, pn);
    $: pages = genPageslist(pn);
    function genPageslist(pn) {
        if (pn < 3) {
            return [1, 2, 3, 4, 5];
        }
        return [pn-2, pn-1, pn, pn+1, pn+2];
    }
    function nextPage() {pn++;}
    function previousPage() {if (pn > 1) pn--;}
    function changePage(p) {pn = p;}

    async function displayPage(q, pn) {
        let response = await fetch("");
        let sharesJSON = await response.json();
        displayResp(sharesJSON);
    }

    function displayResp(json) {
        
    }
</script>

<style>
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
        height: 500px;
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

<div id="main">
    <h1>{name}</h1>
    <form action="submit" autocomplete="on" method="post" accept-charset="utf-8">
        <input type="text"  id="searchbar"/>
        <input type="submit" value="Search" id="searchbutton">
    </form>
        
    <fieldset>
        <legend>Limit</legend>
        <label><input type="radio" name="radio" on:click={() => limit = 10} checked> 10 </label>
        <label><input type="radio" name="radio" on:click={() => limit = 20}> 20 </label>
        <label><input type="radio" name="radio" on:click={() => limit = 50}> 50 </label>
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
        <!-- instantiate other svelte element with right name  -->  
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
            <button on:click={nextPage}>&gt;</button>
        </div>
    {/if}
</div>
