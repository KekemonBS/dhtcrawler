<script>
    export let no;
    export let name;
	export let size;
	export let fileTree;
	export let magnetLink;
</script>

<style>
    div {
        max-width: 100%;
        overflow-x: auto;
        margin-bottom: 1%;
    }

    table {
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
</style>

<div>
    <table>
        <thead>
            <tr>
                <th style="width: 3%">{no}</th>
                <th title={name}>{name}</th>
                <th style="width: 3%">{size}</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td colspan="3" id="#tree">{fileTree}</td>
            </tr>
            <tr>
                <td style="width: 3%;text-align:center;"> 
                    <i class="fa fa-solid fa-magnet"></i>
                </td>
                <td colspan="3">
                    <a href={magnetLink} target="_blank">{magnetLink}</a>
                </td>
            </tr>
        </tbody>
    </table>
</div>
