<script>
    import Post from '../components/Post.vue'
    import { getImgUrl } from '../functions/getImgUrl'
    import { logout } from '../functions/logout'

    export default {
        data() {
            return {
                visitedProfile: {
                    user: {
                        username: '',
                        name: '',
                        birthdate: '',
                        profilePic: '',
                    },
                    posts: []
                },
                newUsername: '',
                loading: true,
                isModifiable: false,
                modifying: false
            }
        },
        methods: {
            getVisitedProfile(){
                this.$axios.get(
                    "/users/" + this.$route.params.username + "/profile/",
                    {
                        headers: {
                            'Authorization': localStorage.getItem('token')
                        },
                    }
                ).then((res) => {
                    this.visitedProfile.user.username = res.data.User.Username
                    this.visitedProfile.user.name = res.data.User.Name
                    this.visitedProfile.user.birthdate = res.data.User.Birthdate
                    this.visitedProfile.user.profilePic = res.data.User.ProfilePic
                    this.visitedProfile.posts = res.data.Posts
                    
                    this.loading = false
                    this.isModifiable = this.visitedProfile.user.username === localStorage.getItem('username')
                }).catch((e) => {
                    alert(e.response.data.ErrorCode + e.response.data.Description)
                })
            },
            updateUsername(){
                this.$axios.put(
                    "/users/" + localStorage.getItem('username') + "/profile/",
                    this.newUsername,
                    {
                        headers: {
                            "Authorization": localStorage.getItem('token'),
                            "Content-Type": 'text/plain'
                        }
                    }
                ).then((res)=>{
                    alert("Username updated! New username: " + this.newUsername)
                    this.newUsername = ''
                    this.loading = true
                    this.modifying = false
                    this.logout()
                }).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
            },
            getImgUrl,
            logout
        },
        created() {  this.getVisitedProfile() },
        updated() { if (!this.loading) {this.getVisitedProfile()} },
        components:{
            Post: Post,
        }
    }
</script>

<template>
    <Header v-if="!loading"></Header>
    <div v-if="!loading"> 
        <span> <img :src="this.getImgUrl(this.visitedProfile.user.profilePic)" style="width: 30px; height: 30px;"> </span>
        <span> 
            <span v-if="!this.modifying"> {{ visitedProfile.user.username }} </span>
            <span v-if="this.modifying">
                <input type="textbox" v-model="newUsername" @keyup.enter="updateUsername">
            </span>
            <img src="@/assets/pencil.jpeg" v-if="isModifiable" @click="this.modifying=!this.modifying" style="width: 30px; height: 30px;">
        </span>
    </div>
    
    <div>
        <Post  
            v-if="visitedProfile.posts"
            v-for="post in visitedProfile.posts" 
            :post="post" 
        ></Post>
    </div>
</template>
