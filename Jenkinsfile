node {
    def root = tool name: 'golang', type: 'go'
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        stage("check out"){
            checkout([$class: 'GitSCM', branches: [[name: '*/dev']], userRemoteConfigs: [[url: 'https://gitee.com/coint/gbbmn-cloud.git', credentialsId: 'GrandpaWang']]])
            // git credentialsId: 'GrandpaWang', url: 'https://gitee.com/coint/gbbmn-cloud.git'
        }

        stage('Build') {
            sh '''
                ls
                cd cmd
                CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloud
                cd ..
            '''
        }

        stage("Deploy in docker"){
            // FIXME(grandpawang)🤦‍jenkins就在服务器的docker 
            // 所以直接去服务器构建和应用
            sh '''
                chmod 777 build/shell/deploy-docker.sh
                ./build/shell/deploy-docker.sh
            '''
        }

        // stage('Deploy') {
        //     sh '''
        //         chmod 777 build/shell/deploy.sh
        //         ./build/shell/deploy.sh
        //     '''
        // }
    }
}