<template>
    <div>
        <div class="tabpanel" v-if="$parent.titleTabActive === 0">
            <mu-data-table
                :columns="logClumns"
                :data="logFiles.slice((currentLogPage-1)*10,currentLogPage*10)"
            >
            <div slot="expand" slot-scope="prop">
                <m-button :href="apiHost+'/logrotate/open?file='+prop.row.Name">打开</m-button>
                <m-button :href="apiHost+'/logrotate/download?file='+prop.row.Name">下载</m-button>
            </div>
                <mu-pagination
                    slot="footer"
                    :total="logFiles.length"
                    :current.sync="currentLogPage"
                ></mu-pagination>
            </mu-data-table>
        </div>
        <div class="tabpanel" v-if="$parent.titleTabActive === 1">
            <div>
                <mu-switch v-model="autoScroll" label="自动滚动" />
            </div>
            <div ref="logContainer" class="log-container">
                <pre><template v-for="item in logs">{{item+"\n"}}</template></pre>
            </div>
        </div>
        <div class="tabpanel" v-if="$parent.titleTabActive === 2">
            <mu-text-field @change="onSearch" placeholder="输入查询关键词"></mu-text-field>
            <pre>{{result}}</pre>
        </div>
    </div>
</template>

<script>
let logsES = null;

export default {
    data() {
        return {
            logClumns:[
                {
                    title:"名称",
                    name:"Name"
                },
                {
                    title:"尺寸",
                    name:"Size",
                    formatter:v=>{
                        return this.unitFormat(v)
                    }
                }
            ],
            autoScroll: true,
            currentLogPage:1,
            logs: [],
            logFiles: [],
            result: "",
            active1: 0
        };
    },
    mounted() {
        logsES = new EventSource(this.apiHost + "/logrotate/tail");
        logsES.onmessage = evt => {
            if (!evt.data) return;
            this.logs.push(evt.data);
        };
        this.ajax
            .getJSON(this.apiHost + "/logrotate/list")
            .then(x => (this.logFiles = x));
        this.$parent.titleTabs=["日志文件","日志跟踪","日志查询"]
    },
    destroyed() {
        logsES.close();
    },
    methods: {
        onSearch(value) {
            this.ajax
                .get(this.apiHost + "/logrotate/find?query=" + value)
                .then(x => (this.result = x));
        }
    },
    updated() {
        if (this.autoScroll && this.$refs.logContainer) {
            this.$refs.logContainer.scrollTop = this.$refs.logContainer.offsetHeight;
        }
    }
};
</script>

<style>
.log-container {
    overflow-y: auto;
    max-height: 500px;
}
.tabpanel {
    padding: 0 20px;
}

.tabpanel .mu-card {
    margin: 5px;
    width: 200px;
}
.tabpanel .mu-card .mu-card-title-container .mu-card-title {
    font-size: 14px;
    line-height: unset;
}
</style>