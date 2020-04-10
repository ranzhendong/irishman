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

        stage('Build') {
            steps {
                echo "3.Build Docker Image Stage"
                sh "docker build -t ${env.HARBOR_URL_TAG}/ranzhendong/irishman:${build_tag} ."
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
             -d '{"msgtype": "text", 
                    "text": {
                    "content": "部署成功"
                    }
                }'
            """
        }
        failure {
            sh """
             curl '${env.DINGTALK_ROBOT}' \
             -H 'Content-Type: application/json' \
             -d '{"msgtype": "text", 
                    "text": {
                    "content": "部署失败"
                    }
                }'
            """
        }
    }
}
