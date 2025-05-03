import joblib
import pandas as pd
import os
import sys
import json

model_dict = joblib.load('scripts/final_model_xgboost.pkl')
scaler = joblib.load('scripts/scaler.pkl')
ordered_features = ['umur', 'monthly_income', 'payroll', 'income_per_age', 'young_rich_flag', 'num_products_owned', 'has_multiple_products', 'income_bucket', 'age_bucket', 'transaction_activity_num', 'income_x_activity', 'gender_MALE', 'marital_status_Single', 'transaction_activity_Inactive', 'categorysegmen_BUMN', 'categorysegmen_Lembaga Negara', 'categorysegmen_Non Target Market', 'categorysegmen_Pendidikan', 'categorysegmen_Pensiun', 'categorysegmen_RS', 'categorysegmen_Swasta']
numerical_cols = [
    'umur', 'monthly_income', 'income_per_age', 'income_x_activity',
    'num_products_owned', 'has_multiple_products', 'income_bucket', 'age_bucket'
]
produk_list = ['griya', 'oto', 'mitraguna', 'hasanahcard', 'pensiun', 'prapensiun']

def predict_final_deploy(data_user):
    input_dict = {
        'umur': data_user['umur'],
        'monthly_income': data_user['income'],
        'payroll': data_user['payroll'],
        'gender_MALE': 1 if data_user['gender'].lower() == 'male' else 0,
        'marital_status_Single': 1 if data_user['marital_status'] == 0 else 0,
        'transaction_activity_Inactive': 1 if data_user['transaction_activity'].lower() == 'inactive' else 0,
        'transaction_activity_num': 0,  # default
    }

    segmen_options = ['BUMN', 'Lembaga Negara', 'Non Target Market', 'Pendidikan', 'Pensiun', 'RS', 'Swasta']
    for segmen in segmen_options:
        col = f'categorysegmen_{segmen}'
        input_dict[col] = 1 if data_user['category_segmen'] == segmen else 0

    # Buat DataFrame
    df = pd.DataFrame([input_dict])

    # Fitur turunan
    df['income_per_age'] = df['monthly_income'] / (df['umur'] + 1)
    df['young_rich_flag'] = ((df['umur'] < 30) & (df['monthly_income'] > 10_000_000)).astype(int)
    df['income_x_activity'] = df['monthly_income'] * df['transaction_activity_num']

    existing_products = data_user['existing_product']
    if not isinstance(existing_products, list):
        existing_products = [existing_products]

    df['num_products_owned'] = len(existing_products)
    df['has_multiple_products'] = int(df['num_products_owned'].values[0] > 1)
    df['income_bucket'] = pd.qcut([df['monthly_income'].values[0]], q=4, labels=False, duplicates='drop')[0]
    df['age_bucket'] = pd.qcut([df['umur'].values[0]], q=4, labels=False, duplicates='drop')[0]
    transaction_map = {'Inactive': 0, 'Active': 1}
    df['transaction_activity_num'] = transaction_map.get(data_user['transaction_activity'], 0)

    for col in ordered_features:
        if col not in df.columns:
            df[col] = 0

    df = df[ordered_features]
    df[numerical_cols] = scaler.transform(df[numerical_cols])

    proba_dict = {}
    for produk in produk_list:
        model = model_dict[produk]
        proba = model.predict_proba(df)[:, 1][0]
        proba_dict[produk] = proba

    if isinstance(data_user['existing_product'], str):
        existing = [data_user['existing_product']]
    else:
        existing = data_user['existing_product']
    for produk in existing:
        if produk in proba_dict:
            proba_dict[produk] = 0
    if data_user['payroll'] == 0 and 'mitraguna' in proba_dict:
        proba_dict['mitraguna'] = 0
    return pd.DataFrame([proba_dict])


if __name__ == "__main__":
    input_json = sys.stdin.read()
    dummy_data = json.loads(input_json)

    dummy_pred_df = predict_final_deploy(dummy_data)
    top_preds = dummy_pred_df.T.sort_values(by=0, ascending=False).head(3)[0].to_dict()
    # print(top_preds)

    print(json.dumps(top_preds))