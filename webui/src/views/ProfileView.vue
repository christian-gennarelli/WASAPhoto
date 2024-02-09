<script>
    import Post from '../components/Post.vue'
    import PopupUserlist from '../components/PopupUserlist.vue'
    import { getImgUrl } from '../functions/getImgUrl'
    import { logout } from '../functions/logout'

    export default {
        data() {
            return {

                // Visited user profile
                visitedProfile: {
                    user: {
                        Username: '',
                        Name: '',
                        Birthdate: '',
                        ProfilePic: '',
                    },
                    posts: [],
                    followers: [],
                    followings: [],
                    banned: [],
                },

                // Authenticated user profile
                authProfile: { // User info
                    user: {
                        ID: localStorage.getItem("ID"),
                        Username: localStorage.getItem("Username"),
                        Name: localStorage.getItem("Name"),
                        Birthdate: localStorage.getItem("Birthdate"),
                        ProfilePic: localStorage.getItem("ProfilePic"),
                    },
                    posts: [],
                    followers: [],
                    followings: [],
                    banned: [],
                },

                // Miscellaneuous
                newUsername: '',
                loading: true,
                modifying: false,
                followed: false,
                showFollowers: false,
                showFollowings: false,
                showBanned: false,
                showUpload: false,
                description: '',
                photo: '',
            }
        },
        computed: {
            isAuthProfile() {
                return this.$route.params.username === this.authProfile.user.Username
            }
        },
        methods: {
            updateUsername(){
                this.$axios.put(
                    "/users/" + this.authProfile.user.Username + "/profile/",
                    this.newUsername,
                    {
                        headers: {
                            "Authorization": this.authProfile.user.ID,
                            "Content-Type": 'text/plain'
                        }
                    }
                ).then(()=>{
                    alert("Username updated! New username: " + this.newUsername)
                    this.newUsername = ''
                    this.loading = true
                    this.modifying = false
                    this.logout()
                }).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
            },
            followUser(){
                this.$axios.put(
                    '/users/' + this.authProfile.user.Username + '/followings/',
                    null,
                    {
                        params: {
                            followed_username: this.visitedProfile.user.Username,
                        },
                        headers: {
                            'Authorization': this.authProfile.user.ID
                        }
                    }
                ).catch(()=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.visitedProfile.followers.unshift(this.authProfile.user)
                this.authProfile.followings.unshift(this.visitedProfile.user)
                this.followed = !this.followed 
            },
            unfollowUser(){
                this.$axios.delete(
                    '/users/' + this.authProfile.user.Username + '/followings/' + this.visitedProfile.user.Username,
                    {
                        headers: {
                            "Authorization": this.authProfile.user.ID,
                        }
                    }
                ).catch(()=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.visitedProfile.followers = this.visitedProfile.followers.filter(u => u.Username !== this.authProfile.user.Username)
                this.authProfile.followings = this.authProfile.followings.filter(u => u.Username !== this.visitedProfile.user.Username)
                this.followed = !this.followed 
            },
            banUser(){
                this.$axios.put(
                    '/users/' + this.authProfile.user.Username + '/banned/',
                    null,
                    {
                        params: {
                            banned_username: this.visitedProfile.user.Username
                        },
                        headers: {
                            'Authorization': this.authProfile.user.ID
                        }
                    }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                alert("You banned this user successfully! Press OK to go back on the home page...")
                this.$router.push('/home')
            },
            unbanUser(bannedUsername){
                this.$axios.delete(
                    '/users/' + this.authProfile.user.Username + '/banned/' + bannedUsername,
                    {
                        headers: {
                            'Authorization': this.authProfile.user.ID
                        }
                    }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.visitedProfile.banned = this.visitedProfile.banned.filter(u => u.Username !== bannedUsername)
            },
            uploadPhoto(){
                const formData = new FormData();
                formData.append('photo', this.photo);
                formData.append('description', this.description)
                this.$axios.post(
                    '/users/' + this.authProfile.user.Username + '/profile/posts/',
                    formData,
                    {
                        headers: {
                            'Authorization': this.authProfile.user.ID,
                            'Content-Type': 'multipart/form-data'
                        }
                    }
                ).then((res)=>{
                    this.visitedProfile.posts ? this.visitedProfile.posts.unshift(res.data) : this.visitedProfile.posts = [res.data]
                }).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.showUpload = false
                this.description = ''
            },
            async getVisitedProfile(){

                // Get authenticated user profile
                await this.$axios.get(
                    "/users/" + this.authProfile.user.Username + "/profile/",
                    {
                        headers: {
                            'Authorization': this.authProfile.user.ID
                        },
                    }
                ).then((res) => {
                    this.authProfile.followings = res.data.Followings ? res.data.Followings : []
                    this.authProfile.followers = res.data.Followers ? res.data.Followers : []
                    this.authProfile.banned = res.data.Banned ? res.data.Banned : []
                    this.authProfile.posts = res.data.Posts ? res.data.Posts : []
                }).catch((e) => {
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                    this.$router.push('/')
                })

                // Get visited user profile
                await this.$axios.get(
                    "/users/" + this.$route.params.username + "/profile/",
                    {
                        headers: {
                            'Authorization': this.authProfile.user.ID
                        },
                    }
                ).then((res) => {
                    this.visitedProfile.user = res.data.User
                    this.visitedProfile.followings = res.data.Followings ? res.data.Followings : []
                    this.visitedProfile.followers = res.data.Followers ? res.data.Followers : []
                    this.visitedProfile.banned = res.data.Banned ? res.data.Banned : []
                    this.visitedProfile.posts = res.data.Posts ? res.data.Posts : []
                    for (let i = 0; i < this.authProfile.followings.length; i++){
                        if (this.authProfile.followings[i].Username == this.visitedProfile.user.Username){
                            this.followed = true
                            break
                        }
                    }
                }).catch((e) => {
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                    this.$router.push('/')
                })
                this.loading = false
            },
            getImgUrl,
            logout
        },
        watch:{
            $route (to, from){
                if (to.name == "profile") {
                    this.loading = true
                    this.followed = false
                    this.getVisitedProfile()
                }
            }
        },
        created() { 
            console.log(this.authProfile.followings)
            this.getVisitedProfile()
        },
        components:{
            Post: Post, PopupUserlist
        }
    }
</script>

<template>
    <Header></Header>
    <div v-if="!loading" class="profile-container"> 
        <div class="left">
            <img class="profile-img" :src="this.getImgUrl(this.visitedProfile.user.ProfilePic)">
            <span class="profile-username" v-if="!this.modifying"> {{ visitedProfile.user.Username }} </span>
            <span class="modified-profile-username" v-else>
                <input type="textbox" v-model="newUsername" @keyup.enter="updateUsername" placeholder="Enter new username">
                <!-- <span> Remember: username can contain only lower/upper case letters, underscores (_) and dashes (-). </span> -->
            </span>
            <img class="pencil" src="@/assets/buttons/pencil.png" v-if="this.isAuthProfile" type="button" @click="this.modifying=!this.modifying">
            <button v-if="followed" class="left-btn btn-red" button @click="unfollowUser"> Unfollow </button>
            <button v-else-if="!this.isAuthProfile" class="left-btn btn-green" @click="followUser"> Follow </button>
            <button v-if="!this.isAuthProfile" class="left-btn btn-red" @click="banUser"> Ban </button>
        </div>

        <div class="right" :style="[!isAuthProfile ? {'grid-template-columns': '1fr 1fr 1fr'} : {'': ''} ]"> 
            <PopupUserlist
                category="Followers"
                :show="this.showFollowers"
                :list="this.visitedProfile.followers"
                @change-show="this.showFollowers=!this.showFollowers; this.showFollowings=false; this.showBanned=false"
            >
            </PopupUserlist>
            <PopupUserlist
                category="Followings"
                :show="this.showFollowings"
                :list="this.visitedProfile.followings"
                @change-show="this.showFollowings=!this.showFollowings; this.showFollowers=false; this.showBanned=false"
            >
            </PopupUserlist>
            <PopupUserlist
                category="Banned"
                :show="this.showBanned"
                :list="this.visitedProfile.banned"
                @change-show="this.showBanned=!this.showBanned; this.showFollowers=false; this.showFollowings=false"
                @unban-user="(bannedUsername) => unbanUser(bannedUsername)"
                @remove-comment="(commentID) => this.visitedProfile.posts[key].Comments = this.visitedProfile.posts[key].Comments.filter(c => c.CommentID != commentID)"
            >
            </PopupUserlist>
            <button v-if="this.visitedProfile.user.Username == this.authProfile.user.Username" @click="this.showUpload=true"> Upload new post </button>
            <div v-if="showUpload" class="upload-overlay">
                <div class="upload-popup">
                    <img type="button" @click="this.showUpload=false" class="exit" src="@/assets/buttons/close.png">
                    <span style="display: block; font-size: 20px; font-weight: bold;"> Photo </span>
                    <input type="file" name="img" accept=".png, .jpeg" ref="file" @change="this.photo=this.$refs.file.files[0]">
                    <span style="display: block; font-size: 20px; font-weight: bold;"> Description </span>
                    <textarea v-model="description" style="display: block; width: 300px; height: 100px" maxlength="128"></textarea>
                    <button style="margin-top: 1%" @click="uploadPhoto"> Upload! </button>
                </div>
            </div>
        </div>
    </div>
    
    <div class="home-container" v-if="!showFollowers && !showFollowings && !showBanned && visitedProfile.posts">
        <Post  
            v-for="(post, key) in visitedProfile.posts" 
            :post="post"
            :user="authProfile.user"
            @add-like="this.visitedProfile.posts[key].Likes ? this.visitedProfile.posts[key].Likes.unshift(this.authProfile.user) : this.visitedProfile.posts[key].Likes = [this.authProfile.user]"
            @remove-like="this.visitedProfile.posts[key].Likes = this.visitedProfile.posts[key].Likes.filter(u => u.Username !== this.authProfile.user.Username)"
            @remove-post="this.visitedProfile.posts.splice(key, 1)"
            @add-comment="(comment) => this.visitedProfile.posts[key].Comments ? this.visitedProfile.posts[key].Comments.unshift(comment) : this.visitedProfile.posts[key].Comments = [comment]"
            @remove-comment="(commentID) => this.visitedProfile.posts[key].Comments = this.visitedProfile.posts[key].Comments.filter(c => c.CommentID != commentID)"
        ></Post>
    </div>
</template>

<style scoped>

.profile-container {
    margin: 15px;
    border-radius: 15px;
    border: 2px solid black;
    display: grid;
    grid-template-columns: 1fr 2fr;
    height: 120px;
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
}

.left {
    display: flex;
    align-items: center;
}

.left-btn {
    border: 2px solid black;
    background-color: transparent;
    font-size: 30px;
    margin-top: 5px;
    margin-right: 2px;
}

.btn-red:hover {
    background: linear-gradient(108.4deg, rgb(253, 44, 56) 3.3%, rgb(176, 2, 12) 98.4%);
}

.btn-green:hover {
    background: linear-gradient(to right, rgb(182, 244, 146), rgb(51, 139, 147));
}

.left .profile-username {
    font-size: 70px;
}

.left .profile-img {
    width: 70px;
    height: 70px;
    border-radius: 50px;
    margin-right: 15px;
}

.left input {
    height: 60px;
}

.left .pencil {
    width: 50px;
    height: 50px;
    margin-left: 15px;
}

.right {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr;
    grid-template-rows: 1fr;
    align-items: center;
    justify-items: center;
}

.right button {
    background-color: transparent;
    border: 1px solid black;
    box-shadow: 2px 2px 2px;
    width: 83%;
}

.right button:hover {
    color: green
}

.upload-overlay {
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background: black;
    overflow:auto;
    z-index: 1;

}

.upload-popup {

    padding: 20px;
    background: #fff;
    border-radius: 5px;
    width: auto;
    min-width: 400px;
    position: relative;
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
    z-index: 1;

}

</style>