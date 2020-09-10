
let next = 1;

$(function () {
    //クッキー取得
    if (typeof $.cookie('userName') !== 'undefined') {
        let userName = $.cookie('userName');
        let userID = $.cookie('userId');
        let profileImg = $.cookie('image');

        //サイドバー情報組み込み
        $('.image').children('img').attr('src', 'https://storage.googleapis.com/' + profileImg);
        $('.sidebar').find('#userName').text(userName);

        //アカウントedit情報組み込み
        if ($('.account_edit_img') != null) {
            $('.account_edit_img').attr('src', 'https://storage.googleapis.com/' + profileImg);
            $('.content').find('.profile-username').text(userName);
            $('#inputName2').val(userName);
        }
    }

    //
    if ($('#question1').length) {
        $('#question1').css('display', 'block');
        if (!$('#question2').length) {
            $('#next').css('display', 'none');
        }
    }


    //画像アップロード機能初期化
    if ($("#settings").length) {
        bsCustomFileInput.init();
    }

    //入力チェック
    $('#quickForm').validate({
        rules: {
            userName: {
                required: true,
                minlength: 2,
            },
            email: {
                required: true,
                email: true,
            },
            password: {
                required: true,
                minlength: 8
            },
            confirmation_password: {
                required: true,
                minlength: 8,
                equalTo: '#password'
            },
        },
        messages: {
            userName: {
                required: "ユーザー名を入力して下さい",
                minlength: "2文字以上入力して下さい"
            },
            email: {
                required: "メールアドレスを入力して下さい",
                email: "有効なメールアドレスを入力して下さい"
            },
            password: {
                required: "パスワードを入力して下さい",
                minlength: "8文字以上入力して下さい"
            },
            confirmation_password: {
                required: "確認用パスワードを入力して下さい",
                minlength: "8文字以上入力して下さい",
                equalTo: "パスワードが一致していません"
            }
        },
        errorElement: 'span',
        errorPlacement: function (error, element) {
            error.addClass('invalid-feedback');
            element.closest('.form-group').append(error);
        },
        highlight: function (element, errorClass, validClass) {
            $(element).addClass('is-invalid');
        },
        unhighlight: function (element, errorClass, validClass) {
            $(element).removeClass('is-invalid');
        }
    });
});

//学習開始確認モーダル

//学習開始（あとでモーダル処理を付ける）
$('.workbook_learning_start').on('click', function () {
    $(this).find('input[type="hidden"]').attr('name', 'bookId');
    const pageURL = $(this).find('a').attr('href');
    $('form').attr('action', pageURL);
    $('form').submit();
    return false
});

//チェックボックスは後で実装
$("checkbox").on("click", function () {
    $('checkbox').prop('checked', false);  //  全部のチェックを外す
    $(this).prop('checked', true);  //  押したやつだけチェックつける
});

//アカウント削除
function DeleteAccount() {
    window.location.href = '/account_delete';
}

//セクション化
function OpenSection(openBtn, closeBtn, openSection) {
    $('#' + openBtn).css('display', 'none');
    $('#' + closeBtn).css('display', 'block');
    $('#' + openSection).css('display', 'block');
}

//非セクション化
function CloseSection(openBtn, closeBtn, closeSection) {
    $('#' + openBtn).css('display', 'block');
    $('#' + closeBtn).css('display', 'none');
    $('#' + closeSection).css('display', 'none');
}

//設問追加
function CloneElement() {
    let a = $('.collapsed-card').last().clone(true);
    let i = Number(a.find('#questionNumber').attr('name'));
    a.find('#questionNumber').attr('name', i + 1);
    a.find('#questionNumber').val(i + 1);
    let aaa = a.find('#questionNumber').val();
    a.find('#AAA').attr('name', i + 1);
    a.find('#answer1').attr('name', i + 1);
    a.find('#answer2').attr('name', i + 1);
    a.find('#answer3').attr('name', i + 1);
    a.find('#answer4').attr('name', i + 1);
    a.find('#BBB').attr('name', i + 1);
    a.find('#card_title')[0].innerHTML = "問" + (i + 1);
    $('.collapsed-card').last().after(a);
}

//設問削除
function DeleteElement() {
    if ($('.aaa > .card').length != 1) {
        $('.aaa > .card').last().remove();
    }
}

//学習:次へ
function Next() {
    $('#question' + next).css('display', 'none');
    $('#question' + (++next)).css('display', 'block');
    $('#back').css('display', 'block');
    if (!$('#question' + (++next)).length) {
        $('#next').css('display', 'none');
        $('#end').css('display', 'block');
        --next;
    } else {
        --next;
    }
}

//学習:戻る
function Back() {
    $('#question' + next).css('display', 'none');
    $('#question' + (--next)).css('display', 'block');
    if (next == 1) {
        $('#back').css('display', 'none');
        $('#next').css('display', 'block');
    }
}

//あとでモーダル処理を付ける
function End() {
    //各ボタンの表示・非表示設定
    $('#back').css('display', 'none');
    $('#end').css('display', 'none');
    $('#end2').css('display', 'block');

    //すべて設問を表示
    for (i = 1; i <= next; i++) {
        $('#question' + i).css('display', 'block');
    }

    //チェックされた選択肢と回答が合っているか答え合わせ


    //答え合わせのレイアウトに変換する

}

//各問題フォルダーページに飛ばす（あとでモーダル表示処理を付ける）
function End2(){
    window.location.href = '/account_delete';
}