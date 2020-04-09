pipeline {
  agent {
    label slave_label()
  }
  stages {
    stage('Clone') {
      parallel {
        stage('Clone') {
          steps {
            git "${git_irishman_url}"
            echo 'rsync'
            echo imageTag()
            script {
              build_tag = imageTag()
            }

          }
        }

        stage('') {
          steps {
            retry(count: 1) {
              git(changelog: true, branch: 'master', url: 'https://github.com/ranzhendong/irishman.git')
            }

          }
        }

      }
    }

    stage('Build') {
      steps {
        echo '3.Build Docker Image Stage'
        sh "docker build -t ${env.HARBOR_URL_TAG}/ranzhendong/irishman:${build_tag} ."
      }
    }

    stage('Push') {
      steps {
        echo '4.Push Docker Image Stage'
        withCredentials(bindings: [usernamePassword(credentialsId: 'zhendongharbor', passwordVariable: 'zhendongharborPassword', usernameVariable: 'zhendongharborUser')]) {
          sh "docker login -u ${zhendongharborUser} -p ${zhendongharborPassword} ${env.HARBOR_URL}"
          sh "docker push ${env.HARBOR_URL_TAG}/ranzhendong/irishman:${build_tag}"
        }

      }
    }

    stage('YAML') {
      steps {
        echo '5. Change YAML File Stage'
        sh "sed -i 's/<BUILD_TAG>/${build_tag}/' /irishman/irishman-deployment.yaml"
        sh "sed -i 's/<BRANCH_NAME>/${env.BRANCH_NAME}/' /irishman/irishman-deployment.yaml"
      }
    }

    stage('DEPLOY') {
      steps {
        sh 'kubectl apply -f /irishman/irishman-deployment.yaml'
      }
    }

  }
  environment {
    git_irishman_url = 'https://gitlab.ranzhendong.com.cn/ranzhendong/irishman.git'
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
      options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        disableConcurrentBuilds()
        timeout(time: 10, unit: 'MINUTES')
        retry(1)
      }
    }