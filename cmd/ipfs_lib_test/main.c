#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include "libipfs.h"

void path(){
	char res[255] = {0};
	char* ipfsPath = "ipfsPath";
	printf("%s\n",ipfsPath);

	GoInt ret = ipfs_path(ipfsPath,res);

	printf("ipfs_path[%d][%s]\n", ret, res);
}

void init() {
	char res[255] = {0};

	GoInt ret = ipfs_init(res);

	printf("ipfs_init[%d][%s]\n", ret, res);
}

void daemon1() {
	char res[255] = {0};

	ipfs_daemon(res);
	printf("[%d]\n", res);
}

void shutdown(){
	char res[255] = {0};
	GoInt ret = ipfs_shutdown(res);
	printf("ipfs_shutdown[%d][%s]\n", ret, res);
}

void add() {
	printf("add........\n");

	GoString hash;
	hash.p = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn";
	hash.n = strlen(hash.p);

	GoString filepath;
	filepath.p = "/123.png";
	filepath.n = strlen(filepath.p);

	GoString str;
	str.p = "./IMG_7714.JPG.png";
	str.n = strlen(str.p);

	char res[255] = {0};
	GoInt ret = ipfs_add(hash, filepath, str, 3, res);
	printf("ipfs_add[%d][%s]\n", ret, res);

	GoString hash2;
	hash2.p = res;
	hash2.n = strlen(res);

	GoString str2;
	str2.p = "./dir";
	str2.n = strlen(str2.p);

	GoString dirpath;
	dirpath.p = "/aaaagd";
	dirpath.n = strlen(dirpath.p);

	char res2[255] = {0};
	GoInt ret2 = ipfs_add(hash2, dirpath, str2, 3, res2);
	printf("ipfs_add[%d][%s]\n", ret2, res2);

	if(strcmp(res2, "QmSxThnifU7Uss2zyS5Q176Tj6GJMSLR6wH2utenmtcngw") == 0)
	{
		printf("ipfs_add success\n");
	} else {
		printf("ipfs_add fail\n");
	}

}

void delete() {
	printf("delete........\n");

	char res[255] = {0};

	GoString root_hash;
	root_hash.p = "QmSxThnifU7Uss2zyS5Q176Tj6GJMSLR6wH2utenmtcngw";
	root_hash.n = strlen(root_hash.p);

	GoString ipfs_path;
	ipfs_path.p = "/aaaagd/IMG_7714.JPG.png";
	ipfs_path.n = strlen(ipfs_path.p);

	GoInt ret = ipfs_delete(root_hash, ipfs_path, 3, res);
	printf("ret[%d][%s]\n", ret, res);

	if (strcmp(res, "QmXhkNgQqtCZeUR1JK13Af1QPGSyLQ6QsqLFaN9VD3vkMp") == 0 ) {
		printf("ipfs_delete success\n");
	} else {
		printf("ipfs_delete fail\n");
	}
}

void move() {
	printf("move........\n");
	char res[255] = {0};
	GoString root_hash;
	root_hash.p = "QmXhkNgQqtCZeUR1JK13Af1QPGSyLQ6QsqLFaN9VD3vkMp";
	root_hash.n = strlen(root_hash.p);

	GoString ipfs_path_src;
	ipfs_path_src.p = "/aaaagd/dir2";
	ipfs_path_src.n = strlen(ipfs_path_src.p);

	GoString ipfs_path_des;
	ipfs_path_des.p = "/dir_helloworld";
	ipfs_path_des.n = strlen(ipfs_path_des.p);

	GoInt ret = ipfs_move(root_hash, ipfs_path_src, ipfs_path_des, 3, res);
	printf("ret[%d][%s]\n", ret, res);

	if (strcmp(res, "QmYUhi1D5r9z8DZM17i3FcBXFGo3qXLY34KbSR3hgWfW6h") == 0) {
		printf("ipfs_move success\n");
	} else {
		printf("ipfs_move fail\n");
	}
}

void shard() {
	printf("shard........\n");
	GoString object_hash;
	object_hash.p = "QmSub8nJ5RUraQ2fHgTdSgqn1rDWDSeJ3mQQvkD5Pi3Chj";
	object_hash.n = strlen(object_hash.p);

	GoString shard_name;
	shard_name.p = "/sharedir";
	shard_name.n = strlen(shard_name.p);

	char res[255] = {0};

	GoInt ret = ipfs_shard(object_hash, shard_name, 3, res);
	printf("ret[%d][%s]\n", ret, res);

	if (strcmp(res, "QmcJJoqoouBDybRxfGzhprx5v8ZXohL41ckCTY6MuzgUiN") == 0) {
		printf("ipfs_shard success\n");
	} else {
		printf("ipfs_shard fail\n");
	}
}

void get() {
	printf("get........\n");

	GoString shard_hash;
	shard_hash.p = "addr://QmcJJoqoouBDybRxfGzhprx5v8ZXohL41ckCTY6MuzgUiN";
	shard_hash.n = strlen(shard_hash.p);

	GoString os_path;
	os_path.p = "./test_get/";
	os_path.n = strlen(os_path.p);

	GoInt ret = ipfs_get(shard_hash, os_path, 3);
	printf("ret[%d]\n", ret);

	if (ret < 0 ) {
		printf("ipfs_get fail\n");
	} else {
		printf("ipfs_get success\n");
	}
}

void query() {
	printf("query........\n");
	char res[2048] = {0};
	GoString object_hash;
	object_hash.p = "QmYUhi1D5r9z8DZM17i3FcBXFGo3qXLY34KbSR3hgWfW6h";
	object_hash.n = strlen(object_hash.p);

	GoString ipfs_path;
	ipfs_path.p = "/aaaagd";
	ipfs_path.n = strlen(ipfs_path.p);

	GoInt ret = ipfs_query(object_hash, ipfs_path, 3, res);
	printf("ret[%d][%s]\n", ret, res);

	if (strcmp(res, "{\"Objects\":[{\"Hash\":\"QmYrUdKjuKfjN71R2Kt7YQJY5Y3gVUCPLivGbeB3L4zq5t\",\"Links\":[{\"Name\":\"a.txt\",\"Hash\":\"QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o\",\"Size\":20,\"Type\":2},{\"Name\":\"abbbb\",\"Hash\":\"QmSs1VAsZM1Rw4qJeb4fSK8rn69LwMtTwgAeHX96GL8pb3\",\"Size\":29107323,\"Type\":2},{\"Name\":\"b.txt\",\"Hash\":\"QmTgKghvimxUPwVPiTwgiATDQhqUxcWrm1bM6M2cdK7ycM\",\"Size\":15,\"Type\":2},{\"Name\":\"bigdir\",\"Hash\":\"QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn\",\"Size\":4,\"Type\":1},{\"Name\":\"c.txt\",\"Hash\":\"QmU2viJUREiEcFPgx5rzEVtB8psw1F49hYc5C7pS6pffpt\",\"Size\":21,\"Type\":2},{\"Name\":\"dir1\",\"Hash\":\"QmbpFbWtjzaZJtPeWeUaxFoDdDTrEE2wJTNGuDUkJJRZ4P\",\"Size\":1003,\"Type\":1},{\"Name\":\"dir_nil\",\"Hash\":\"QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn\",\"Size\":4,\"Type\":1},{\"Name\":\"libtest\",\"Hash\":\"QmX6r5YamtKTAAVMGyVkxvMPgXQkcCBYvPV7PvR91CMM2F\",\"Size\":19122630,\"Type\":2},{\"Name\":\"nil.txt\",\"Hash\":\"QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH\",\"Size\":6,\"Type\":2}]}]}") == 0 ) {
		printf("ipfs_query success\n");
	} else {
		printf("ipfs_query fail\n");
	}
}

void merge() {
	printf("merge........\n");
	char res[255] = {0};
	GoString root_hash;
	root_hash.p = "QmYUhi1D5r9z8DZM17i3FcBXFGo3qXLY34KbSR3hgWfW6h";
	root_hash.n = strlen(root_hash.p);

	GoString ipfs_path;
	ipfs_path.p = "/dir_helloworld/ttttttt";
	ipfs_path.n = strlen(ipfs_path.p);

	GoString shard_hash;
	shard_hash.p = "QmcJJoqoouBDybRxfGzhprx5v8ZXohL41ckCTY6MuzgUiN";
	shard_hash.n = strlen(shard_hash.p);

	GoInt ret = ipfs_merge(root_hash, ipfs_path, shard_hash, 3, res);
	printf("ret[%d][%s]\n", ret, res);

	if (strcmp(res, "QmRnJKjqNfVN25SrXgbA8KWv4szc5xccWFcx9bwSXvAoT8") == 0) {
		printf("ipfs_merge success\n");
	} else {
		printf("ipfs_merge fail\n");
	}
}

void id() {
	printf("id........\n");
	char res[2048] = {0};
	GoInt ret = ipfs_id(3, res);
	printf("ret[%d][%s]\n", ret, res);
}

void peerid() {
	printf("peerid........\n");
	char res[255] = {0};
	GoString new_hash;
	new_hash.p = "QmWCwooARFCZNhy94MrRwfyhorvwdJWuSHYCTfSrG2VfLg";
	new_hash.n = 0;

	GoInt ret = ipfs_peerid(new_hash, 3, res);
	printf("ret[%d][%s]\n", ret, res);

}

void privkey() {
	printf("privkey........\n");
	char res[2048] = {0};
	GoString new_key;
	new_key.p = "QmWCwooARFCZNhy94MrRwfyhorvwdJWuSHYCTfSrG2VfLg";
	new_key.n = 0;

	GoInt ret = ipfs_privkey(new_key, 3, res);
	printf("ret[%d][%s]\n", ret, res);
}

void publish() {
	printf("publish.......\n");
	GoString hash;
	hash.p = "QmRnJKjqNfVN25SrXgbA8KWv4szc5xccWFcx9bwSXvAoT8";
	hash.n = strlen(hash.p);

	char pRest[255] = {0};
	GoInt ret = ipfs_publish(hash, 0, pRest);
	printf("ret[%d][%s]\n", ret, pRest);

}

void remotepin() {
	printf("remote.........\n");

	GoString peer;
	peer.p = "QmXF8BwQ6BRUrdxgs15NU4h3Pmk8w4duZK3CJdZKuMBQYx";
	peer.n = strlen(peer.p);

	GoString key;
	key.p = "AVj7cRpH";
	key.n = strlen(key.p);

	GoString object;
	object.p = "/ipfs/QmRnJKjqNfVN25SrXgbA8KWv4szc5xccWFcx9bwSXvAoT8";
	object.n = strlen(object.p);

	char pRest[255] = {0};
	GoInt ret = ipfs_remotepin(peer, key, object, 0, pRest);
	printf("ret[%d][%s]success\n", ret, pRest);

	GoString key2;
	key2.p = "DeTh/yo8";
	key2.n = strlen(key2.p);

	char pRest2[255] = {0};
	GoInt ret2 = ipfs_remotepin(peer, key2, object, 0, pRest2);
	printf("ret[%d][%s]fail\n", ret2, pRest2);
}

void remotels() {
	GoString peer;
	peer.p = "QmXF8BwQ6BRUrdxgs15NU4h3Pmk8w4duZK3CJdZKuMBQYx";
	peer.n = strlen(peer.p);

	GoString sec;
	sec.p = "AVj7cRpH";
	sec.n = strlen(sec.p);

	GoString hash;
	hash.p = "QmRnJKjqNfVN25SrXgbA8KWv4szc5xccWFcx9bwSXvAoT8";
	hash.n = strlen(hash.p);

	char pRest2[255] = {0};
	GoInt ret2 = ipfs_remotels(peer, sec, hash, 0, pRest2);
	printf("ret[%d][%s]\n", ret2, pRest2);
}

void connect() {
	GoString peer;
	peer.p = "/ip4/172.16.158.1/tcp/4001/ipfs/QmXF8BwQ6BRUrdxgs15NU4h3Pmk8w4duZK3CJdZKuMBQYx";
	peer.n = strlen(peer.p);

	char pRest[255] = {0};
	GoInt ret = ipfs_connectpeer(peer, 0, pRest);
	printf("ret[%d][%s]\n", ret, pRest);
}

void config() {
	printf("config........\n");
	GoString key;
	key.p = "Identity.PeerID";
	key.n = 0;

	GoString value;
	value.p = "QmYboPwvU7wHzdadfZHEfehRSSYx7zjfg8KxgQ1vdQ1pax";
	value.n = 0;

	char res[4096] = {0};

	GoInt ret = ipfs_config(key, value, res);
	printf("ret[%d]res[%s]\n", ret, res);

	memset(res, 0, sizeof(res));
	key.n = strlen(key.p);
	ret = ipfs_config(key, value, res);
	printf("ret[%d]res[%s]\n", ret, res);

	memset(res, 0, sizeof(res));
	value.n = strlen(value.p);
	ret = ipfs_config(key, value, res);
	printf("ret[%d]res[%s]\n", ret, res);
}

void cmd() {
	printf("cmd........\n");
	GoString cmd;
	cmd.p = "ipfs daemon";
	cmd.n = strlen(cmd.p);

	char pRest[255];
	memset(pRest, 0, sizeof(pRest));

	GoInt ret = ipfs_cmd(cmd, 3, pRest);
	printf("ret[%d][%s]\n", ret, pRest);
}

void operate(){
	add();
	delete();
	move();
	shard();
	get();
	query();
	merge();
	peerid();
	privkey();
	// connect();
	// publish();
	// remotepin();
	// remotels();
}

int main() {
	path();
	init();
	
	// do daemon
	pthread_t tid1,tid2;
    int ret1,ret2;
    ret1 = pthread_create(&tid1, NULL, (void*)daemon1, NULL);
    if(ret1 != 0)
    {
        printf("Thread Create Daemon Error\n");
        exit(0);
    }

    // wait for daemon start
    sleep(10);

    ret2 = pthread_create(&tid2, NULL, (void*)operate, NULL);
    if(ret2 != 0)
    {
        printf("Thread Create Operate Error\n");
        exit(0);
    }

	pthread_join(tid2, NULL);
	pthread_join(tid1, NULL);
	shutdown();
}
