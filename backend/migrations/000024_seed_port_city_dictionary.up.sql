INSERT INTO custom_destinations (name, country, latitude, longitude, keywords, description, status, sort_order)
SELECT seed.name, seed.country, seed.latitude, seed.longitude, seed.keywords, seed.description, seed.status, seed.sort_order
FROM (
    SELECT '上海' AS name, '中国' AS country, 31.2304 AS latitude, 121.4737 AS longitude, '上海,shanghai' AS keywords, 'system_port_city_dictionary_seed_v1' AS description, 1 AS status, 1000 AS sort_order
    UNION ALL SELECT '天津', '中国', 39.0842, 117.2009, '天津,tianjin', 'system_port_city_dictionary_seed_v1', 1, 995
    UNION ALL SELECT '大连', '中国', 38.9140, 121.6147, '大连,dalian', 'system_port_city_dictionary_seed_v1', 1, 990
    UNION ALL SELECT '青岛', '中国', 36.0671, 120.3826, '青岛,qingdao', 'system_port_city_dictionary_seed_v1', 1, 985
    UNION ALL SELECT '厦门', '中国', 24.4798, 118.0894, '厦门,xiamen', 'system_port_city_dictionary_seed_v1', 1, 980
    UNION ALL SELECT '香港', '中国', 22.3193, 114.1694, '香港,hong kong,hk', 'system_port_city_dictionary_seed_v1', 1, 975
    UNION ALL SELECT '基隆', '中国', 25.1276, 121.7392, '基隆,keelung', 'system_port_city_dictionary_seed_v1', 1, 970
    UNION ALL SELECT '高雄', '中国', 22.6273, 120.3014, '高雄,kaohsiung', 'system_port_city_dictionary_seed_v1', 1, 965
    UNION ALL SELECT '福冈', '日本', 33.5902, 130.4017, '福冈,博多,fukuoka,hakata', 'system_port_city_dictionary_seed_v1', 1, 960
    UNION ALL SELECT '长崎', '日本', 32.7503, 129.8777, '长崎,nagasaki', 'system_port_city_dictionary_seed_v1', 1, 955
    UNION ALL SELECT '鹿儿岛', '日本', 31.5966, 130.5571, '鹿儿岛,kagoshima', 'system_port_city_dictionary_seed_v1', 1, 950
    UNION ALL SELECT '横滨', '日本', 35.4437, 139.6380, '横滨,yokohama', 'system_port_city_dictionary_seed_v1', 1, 945
    UNION ALL SELECT '神户', '日本', 34.6901, 135.1955, '神户,kobe', 'system_port_city_dictionary_seed_v1', 1, 940
    UNION ALL SELECT '大阪', '日本', 34.6937, 135.5023, '大阪,osaka', 'system_port_city_dictionary_seed_v1', 1, 935
    UNION ALL SELECT '济州', '韩国', 33.4996, 126.5312, '济州,济州岛,jeju', 'system_port_city_dictionary_seed_v1', 1, 930
    UNION ALL SELECT '西归浦', '韩国', 33.2541, 126.5601, '西归浦,西归浦市,seogwipo', 'system_port_city_dictionary_seed_v1', 1, 925
    UNION ALL SELECT '釜山', '韩国', 35.1796, 129.0756, '釜山,busan', 'system_port_city_dictionary_seed_v1', 1, 920
    UNION ALL SELECT '仁川', '韩国', 37.4563, 126.7052, '仁川,incheon', 'system_port_city_dictionary_seed_v1', 1, 915
    UNION ALL SELECT '新加坡', '新加坡', 1.2903, 103.8519, '新加坡,singapore', 'system_port_city_dictionary_seed_v1', 1, 910
    UNION ALL SELECT '槟城', '马来西亚', 5.4164, 100.3327, '槟城,penang', 'system_port_city_dictionary_seed_v1', 1, 905
    UNION ALL SELECT '巴生港', '马来西亚', 3.0019, 101.3910, '巴生港,port klang', 'system_port_city_dictionary_seed_v1', 1, 900
    UNION ALL SELECT '普吉', '泰国', 7.8804, 98.3923, '普吉,phuket', 'system_port_city_dictionary_seed_v1', 1, 895
    UNION ALL SELECT '林查班', '泰国', 13.0827, 100.8830, '林查班,laem chabang', 'system_port_city_dictionary_seed_v1', 1, 890
    UNION ALL SELECT '迈阿密', '美国', 25.7617, -80.1918, '迈阿密,邁阿密,miami', 'system_port_city_dictionary_seed_v1', 1, 885
    UNION ALL SELECT '劳德代尔堡', '美国', 26.1224, -80.1373, '劳德代尔堡,罗德岱堡,fort lauderdale,port everglades', 'system_port_city_dictionary_seed_v1', 1, 880
    UNION ALL SELECT '卡纳维拉尔港', '美国', 28.4089, -80.6043, '卡纳维拉尔港,奥兰多港,port canaveral', 'system_port_city_dictionary_seed_v1', 1, 875
    UNION ALL SELECT '拿骚', '巴哈马', 25.0443, -77.3504, '拿骚,nassau', 'system_port_city_dictionary_seed_v1', 1, 870
    UNION ALL SELECT '科苏梅尔', '墨西哥', 20.4229839, -86.9223432, '科苏梅尔,cozumel,isla cozumel', 'system_port_city_dictionary_seed_v1', 1, 865
    UNION ALL SELECT '乔治城', '开曼群岛', 19.2866, -81.3674, '乔治城,george town,grand cayman', 'system_port_city_dictionary_seed_v1', 1, 860
    UNION ALL SELECT '圣胡安', '波多黎各', 18.4655, -66.1057, '圣胡安,san juan', 'system_port_city_dictionary_seed_v1', 1, 855
    UNION ALL SELECT '巴塞罗那', '西班牙', 41.3851, 2.1734, '巴塞罗那,barcelona', 'system_port_city_dictionary_seed_v1', 1, 850
    UNION ALL SELECT '奇维塔韦基亚', '意大利', 42.0924, 11.7950, '奇维塔韦基亚,罗马港,civitavecchia', 'system_port_city_dictionary_seed_v1', 1, 845
    UNION ALL SELECT '那不勒斯', '意大利', 40.8518, 14.2681, '那不勒斯,naples,napoli', 'system_port_city_dictionary_seed_v1', 1, 840
    UNION ALL SELECT '比雷埃夫斯', '希腊', 37.9420, 23.6465, '比雷埃夫斯,雅典港,piraeus', 'system_port_city_dictionary_seed_v1', 1, 835
    UNION ALL SELECT '马赛', '法国', 43.2965, 5.3698, '马赛,marseille', 'system_port_city_dictionary_seed_v1', 1, 830
    UNION ALL SELECT '里斯本', '葡萄牙', 38.7223, -9.1393, '里斯本,lisbon,lisboa', 'system_port_city_dictionary_seed_v1', 1, 825
    UNION ALL SELECT '南安普敦', '英国', 50.9097, -1.4044, '南安普敦,southampton', 'system_port_city_dictionary_seed_v1', 1, 820
    UNION ALL SELECT '布宜诺斯艾利斯', '阿根廷', -34.6037, -58.3816, '布宜诺斯艾利斯,布宜諾斯艾利斯,buenos aires', 'system_port_city_dictionary_seed_v1', 1, 815
    UNION ALL SELECT '蒙得维的亚', '乌拉圭', -34.9011, -56.1645, '蒙得维的亚,montevideo', 'system_port_city_dictionary_seed_v1', 1, 810
    UNION ALL SELECT '桑托斯', '巴西', -23.9608, -46.3336, '桑托斯,santos', 'system_port_city_dictionary_seed_v1', 1, 805
    UNION ALL SELECT '里约热内卢', '巴西', -22.9068, -43.1729, '里约热内卢,rio de janeiro', 'system_port_city_dictionary_seed_v1', 1, 800
    UNION ALL SELECT '悉尼', '澳大利亚', -33.8688, 151.2093, '悉尼,sydney', 'system_port_city_dictionary_seed_v1', 1, 795
    UNION ALL SELECT '奥克兰', '新西兰', -36.8509, 174.7645, '奥克兰,auckland', 'system_port_city_dictionary_seed_v1', 1, 790
) AS seed
WHERE NOT EXISTS (
    SELECT 1
    FROM custom_destinations existing
    WHERE existing.name = seed.name
      AND existing.country = seed.country
      AND existing.deleted_at IS NULL
);