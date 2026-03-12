INSERT INTO custom_destinations (name, country, latitude, longitude, keywords, description, status, sort_order)
SELECT seed.name, seed.country, seed.latitude, seed.longitude, seed.keywords, seed.description, seed.status, seed.sort_order
FROM (
    SELECT '温哥华' AS name, '加拿大' AS country, 49.2827 AS latitude, -123.1207 AS longitude, '温哥华,vancouver' AS keywords, 'system_port_city_dictionary_seed_v2' AS description, 1 AS status, 780 AS sort_order
    UNION ALL SELECT '维多利亚', '加拿大', 48.4284, -123.3656, '维多利亚,victoria bc', 'system_port_city_dictionary_seed_v2', 1, 775
    UNION ALL SELECT '朱诺', '美国', 58.3019, -134.4197, '朱诺,juneau', 'system_port_city_dictionary_seed_v2', 1, 770
    UNION ALL SELECT '斯卡格威', '美国', 59.4583, -135.3139, '斯卡格威,skagway', 'system_port_city_dictionary_seed_v2', 1, 765
    UNION ALL SELECT '凯奇坎', '美国', 55.3422, -131.6461, '凯奇坎,ketchikan', 'system_port_city_dictionary_seed_v2', 1, 760
    UNION ALL SELECT '惠蒂尔', '美国', 60.7743, -148.6837, '惠蒂尔,whittier alaska', 'system_port_city_dictionary_seed_v2', 1, 755
    UNION ALL SELECT '西雅图', '美国', 47.6062, -122.3321, '西雅图,seattle', 'system_port_city_dictionary_seed_v2', 1, 750
    UNION ALL SELECT '哥本哈根', '丹麦', 55.6761, 12.5683, '哥本哈根,copenhagen', 'system_port_city_dictionary_seed_v2', 1, 745
    UNION ALL SELECT '斯德哥尔摩', '瑞典', 59.3293, 18.0686, '斯德哥尔摩,stockholm', 'system_port_city_dictionary_seed_v2', 1, 740
    UNION ALL SELECT '赫尔辛基', '芬兰', 60.1699, 24.9384, '赫尔辛基,helsinki', 'system_port_city_dictionary_seed_v2', 1, 735
    UNION ALL SELECT '奥斯陆', '挪威', 59.9139, 10.7522, '奥斯陆,oslo', 'system_port_city_dictionary_seed_v2', 1, 730
    UNION ALL SELECT '雷克雅未克', '冰岛', 64.1466, -21.9426, '雷克雅未克,reykjavik', 'system_port_city_dictionary_seed_v2', 1, 725
    UNION ALL SELECT '塔林', '爱沙尼亚', 59.4370, 24.7536, '塔林,tallinn', 'system_port_city_dictionary_seed_v2', 1, 720
    UNION ALL SELECT '都柏林', '爱尔兰', 53.3498, -6.2603, '都柏林,dublin', 'system_port_city_dictionary_seed_v2', 1, 715
    UNION ALL SELECT '威尼斯', '意大利', 45.4408, 12.3155, '威尼斯,venice,venezia', 'system_port_city_dictionary_seed_v2', 1, 710
    UNION ALL SELECT '热那亚', '意大利', 44.4056, 8.9463, '热那亚,genoa,genova', 'system_port_city_dictionary_seed_v2', 1, 705
    UNION ALL SELECT '帕尔马', '西班牙', 39.5696, 2.6502, '帕尔马,palma de mallorca', 'system_port_city_dictionary_seed_v2', 1, 700
    UNION ALL SELECT '瓦莱塔', '马耳他', 35.8989, 14.5146, '瓦莱塔,valletta', 'system_port_city_dictionary_seed_v2', 1, 695
    UNION ALL SELECT '圣托里尼', '希腊', 36.3932, 25.4615, '圣托里尼,santorini,thira', 'system_port_city_dictionary_seed_v2', 1, 690
    UNION ALL SELECT '米科诺斯', '希腊', 37.4467, 25.3289, '米科诺斯,mykonos', 'system_port_city_dictionary_seed_v2', 1, 685
    UNION ALL SELECT '伊斯坦布尔', '土耳其', 41.0082, 28.9784, '伊斯坦布尔,istanbul', 'system_port_city_dictionary_seed_v2', 1, 680
    UNION ALL SELECT '多哈', '卡塔尔', 25.2854, 51.5310, '多哈,doha', 'system_port_city_dictionary_seed_v2', 1, 675
    UNION ALL SELECT '迪拜', '阿联酋', 25.2048, 55.2708, '迪拜,dubai', 'system_port_city_dictionary_seed_v2', 1, 670
    UNION ALL SELECT '阿布扎比', '阿联酋', 24.4539, 54.3773, '阿布扎比,abu dhabi', 'system_port_city_dictionary_seed_v2', 1, 665
    UNION ALL SELECT '布里奇顿', '巴巴多斯', 13.0975, -59.6167, '布里奇顿,bridgetown', 'system_port_city_dictionary_seed_v2', 1, 660
    UNION ALL SELECT '菲利普斯堡', '荷属圣马丁', 18.0260, -63.0458, '菲利普斯堡,philipsburg,st maarten', 'system_port_city_dictionary_seed_v2', 1, 655
    UNION ALL SELECT '法尔茅斯', '牙买加', 18.4928, -77.6563, '法尔茅斯,falmouth jamaica', 'system_port_city_dictionary_seed_v2', 1, 650
    UNION ALL SELECT '罗阿坦', '洪都拉斯', 16.3170, -86.5371, '罗阿坦,roatan', 'system_port_city_dictionary_seed_v2', 1, 645
) AS seed
WHERE NOT EXISTS (
    SELECT 1
    FROM custom_destinations existing
    WHERE existing.name = seed.name
      AND existing.country = seed.country
      AND existing.deleted_at IS NULL
);