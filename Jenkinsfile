node ('jenkins-slave-k8s'){

    stage('ECHO') {
        sh " echo ${env.HARBOR_URL}"
        sh "echo '192.168.10.10   ${env.HARBOR_URL_TAG}' >>/etc/hosts "
        sh "cat /etc/hosts"
    }

    stage('Test') {
        echo "2.Test Stage"
    }

    stage('Build') {
        echo "3.Build Docker Image Stage"


        sh "docker build -t ${env.HARBOR_URL_TAG}/${env.IRISHMAN_HARBOR_IMAGE}:${build_tag} ."
    }

    stage('Push') {
        echo "4.Push Docker Image Stage"
        withCredentials(
            [usernamePassword(credentialsId: 'zhendongharbor', passwordVariable: 'zhendongharborPassword', usernameVariable: 'zhendongharborUser')]) {
            sh "docker login -u ${zhendongharborUser} -p ${zhendongharborPassword} ${env.HARBOR_URL}"
            sh "docker push ${env.HARBOR_URL_TAG}/${env.IRISHMAN_HARBOR_IMAGE}:${build_tag}"
        }
    }

    stage('YAML') {
        echo "5. Change YAML File Stage"

        sh "sed -i 's!image:.*!image: ${env.HARBOR_URL_TAG}/${env.IRISHMAN_HARBOR_IMAGE}:${build_tag}!' /irishman/irishman-deployment.yaml"
        sh "sed -i 's/value: .*/value: ${env.BRANCH_NAME}/' /irishman/irishman-deployment.yaml"

        sh "kubectl apply -f /irishman/irishman-deployment.yaml"
    }

}