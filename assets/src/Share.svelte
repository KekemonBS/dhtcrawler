<script>
    export let no;
    export let Name;
	export let Size;
	export let FileTree;
	export let MagnetLink;
</script>

<style>
    div {
        max-width: 100%;
        overflow-x: auto;
        margin-bottom: 1%;
    }

    table {
        min-width: 100%;
        background-color: #eeeeee;
        border: 2px solid;
        border-collapse: collapse;
    }

    tr,th,td {
        border: 1px solid;
        border-collapse: collapse;
        padding: 3pt;
    }

    th {
        text-align: center;
        background-color: #dedede;
    }

    #filetree {
        white-space: pre;
        max-height: 300px; 
        overflow-y: auto;
        overflow-x: hidden;
    }
</style>

<div>
    <table>
        <thead>
            <tr>
                <th style="width: 3%">{no}</th>
                <th title={Name}>{Name}</th>
                <th style="width: 3%">{Size}</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td colspan="3" id="#tree">
                    <div id=filetree>
                        {FileTree}
                    </div>
                </td>
            </tr>
            <tr>
                <td style="width: 3%;text-align:center;"> 
                    <i class="fa fa-solid fa-magnet"></i>
                </td>
                <td colspan="3">
                    <a href={MagnetLink} target="_blank">{MagnetLink}</a>
                </td>
            </tr>
        </tbody>
    </table>
</div>
