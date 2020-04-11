def slave_label() {
    return "jenkins-slave-k8s"
    }

def imageTag() {
    return  sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
    }

def createVersion() {
    // 定义一个版本号作为当次构建的版本，输出结果 20191210175842_69
    return new Date().format('yyyyMMddHHmmss') + "_${env.BUILD_ID}"
    }

def buildImage() {
    // docker.build("${env.HARBOR_URL_TAG}/ranzhendong/irishman:${build_tag}")
    // script表示里面是脚本式写法
        script {
            docker.build("${env.HARBOR_URL_TAG}/ranzhendong/irishman:${build_tag}")
        }
    }


pipeline {
    // agent any
    agent{
        label slave_label()
    }

    options {
         // 表示保留10次构建历史
        buildDiscarder(logRotator(numToKeepStr: '10'))

        // 不允许同时执行流水线，被用来防止同时访问共享资源等
        disableConcurrentBuilds()

        // 设置流水线运行的超时时间, 在此之后，Jenkins将中止流水线
        timeout(time: 10, unit: 'MINUTES')

        // 重试次数
        retry(1)
        
    }
    
    stages {

        stage('Add ENV'){
            steps{
                sh 'printenv'
                // 脚本式写法，赋值变量
                script {
                    build_tag = imageTag()
                    //dingtalk robot infomation
                    title = "Jenkins Notice"
                    successtext = """
### 【${env.gitlabSourceRepoName} 构建Success】\n\n
#### 构建人：${env.gitlabUserName}\n
> 触发分支：${env.gitlabTargetBranch}\n
> 项目地址：[${env.gitlabSourceRepoName}](${env.gitlabSourceRepoHomepage})\n
> COMMIT地址：[${env.gitlabMergeRequestLastCommit}](${env.gitlabSourceRepoHomepage}/commit/${env.gitlabMergeRequestLastCommit})\n
"""
                    failuretext = """
### 【${env.gitlabSourceRepoName} 构建Success】\n\n
#### 构建人：${env.gitlabUserName}\n
> 触发分支：${env.gitlabTargetBranch}\n
> 项目地址：[${env.gitlabSourceRepoName}](${env.gitlabSourceRepoHomepage})\n
> COMMIT地址：[${env.gitlabMergeRequestLastCommit}](${env.gitlabSourceRepoHomepage}/commit/${env.gitlabMergeRequestLastCommit})\n
"""
                }
            }
        }

        stage('Build') {
            steps {
                // 镜像构建
                buildImage()
            }
        }
        
        stage('Push') {
            steps {
                echo "4.Push Docker Image Stage"
                withCredentials([usernamePassword(credentialsId: 'zhendongharbor', passwordVariable: 'zhendongharborPassword', usernameVariable: 'zhendongharborUser')]) {
                sh "docker login -u ${zhendongharborUser} -p ${zhendongharborPassword} ${env.HARBOR_URL}"
                sh "docker push ${env.HARBOR_URL_TAG}/ranzhendong/irishman:${build_tag}"
                }
            }
        }
        
        stage('YAML') {
            steps {
                echo "5. Change YAML File Stage"
                sh "sed -i 's/<BUILD_TAG>/${build_tag}/' /irishman/irishman-deployment.yaml"
                sh "sed -i 's/<BRANCH_NAME>/${env.BRANCH_NAME}/' /irishman/irishman-deployment.yaml"
            }
        }
        
        stage('DEPLOY') {
            steps {
                sh "kubectl apply -f /irishman/irishman-deployment.yaml"
            }
        }
        
    }
    post {
        success {
            sh """
             curl '${env.DINGTALK_ROBOT}' \
             -H 'Content-Type: application/json' \
             -d '{"msgtype": "markdown", 
                "markdown": {
                    "title": "${title}",
                    "text": "${successtext}"
                    }
                }'
            """
        }
        failure {
            sh """
             curl '${env.DINGTALK_ROBOT}' \
             -H 'Content-Type: application/json' \
             -d '{"msgtype": "markdown", 
                "markdown": {
                    "title": "${title}",
                    "text": "${failuretext}"
                    }
                }'
            """
        }
    }
}
