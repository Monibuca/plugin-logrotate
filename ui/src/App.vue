<template>
    <div>
        <mu-tabs :value.sync="active1" indicator-color="#80deea" inverse center>
            <mu-tab>日志文件</mu-tab>
            <mu-tab>日志跟踪</mu-tab>
            <mu-tab>日志查询</mu-tab>
        </mu-tabs>
        <div class="tabpanel tab1" v-if="active1 === 0">
            <mu-card v-for="item in logFiles" :key="item.Name">
                <mu-card-title :title="item.Name" :sub-title="unitFormat(item.Size)"></mu-card-title>
                <mu-card-actions>
                    <mu-button small flat :href="apiHost+'/logrotate/open?file='+item.Name" target="_blank">打开
                    </mu-button>
                    <mu-button small flat :href="apiHost+'/logrotate/download?file='+item.Name" target="_blank">下载
                    </mu-button>
                </mu-card-actions>
            </mu-card>
        </div>
        <div class="tabpanel" v-if="active1 === 1">
            <div>
                <mu-switch v-model="autoScroll" label="自动滚动" />
            </div>
            <div ref="logContainer" class="log-container">
                <pre><template v-for="item in logs">{{item+"\n"}}</template></pre>
            </div>
        </div>
        <div class="tabpanel" v-if="active1 === 2">
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
            autoScroll: true,
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
.tab1 {
    display: flex;
    flex-wrap: wrap;
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